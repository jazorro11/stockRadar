// Declaramos que este archivo pertenece al paquete principal (main), obligatorio para ejecutar un programa Go.
package main

// Importamos los paquetes necesarios para hacer solicitudes HTTP, leer datos, manejar errores, parsear JSON, etc.
import (
	"encoding/json" // Permite convertir entre JSON y estructuras Go.
	"fmt"           // Sirve para imprimir texto en consola.
	"io"            // Proporciona utilidades para leer el cuerpo de la respuesta HTTP.
	"net/http"      // Permite crear clientes HTTP y hacer solicitudes.
	"os"            // Permite interactuar con el sistema operativo (ej: salir del programa con os.Exit).
	"strconv"       // Permite convertir cadenas a otros tipos de datos.
	"strings"       // Permite manipular cadenas de texto.

	"github.com/joho/godotenv" // opcional, para cargar .env en local
	// This will create a go.mod file and allow you to use external packages like github.com/joho/godotenv.
	//go mod init stockradar
	//go get github.com/joho/godotenv
	"context"

	"github.com/jackc/pgx/v5" // Importa el paquete pgx para manejar conexiones a bases de datos PostgreSQL.
	//go get github.com/jackc/pgx/v5
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

// Definimos la estructura que representa la respuesta principal de la API.
// Esta estructura espera un campo "items" que contiene un arreglo de acciones (StockInfo).
type ApiResponse struct {
	Items    []StockInfo `json:"items"`      // Lista de acciones recibidas del API.
	NextPage string      `json:"next_page"`  // Paginación, por si hay más datos (opcional).
}

// Definimos la estructura que representa cada acción (stock).
// Los tags `json:"nombre"` permiten mapear los campos del JSON con los campos de la estructura Go.
type StockInfo struct {
	Ticker      string `json:"ticker"`       // Símbolo bursátil de la acción.
	Company     string `json:"company"`      // Nombre de la empresa.
	Brokerage   string `json:"brokerage"`    // Nombre del bróker que reporta la recomendación.
	StockAction string `json:"action"`       // Tipo de recomendación (ej: "Buy", "Hold", etc).
	RatingFrom  string `json:"rating_from"`  // Calificación anterior del bróker.
	RatingTo    string `json:"rating_to"`    // Nueva calificación otorgada.
	TargetFrom  float64 `json:"target_from"`  // Precio objetivo anterior.
	TargetTo    float64 `json:"target_to"`    // Precio objetivo nuevo.
	StockTime   string `json:"time"`   // Momento en que se emitió la recomendación.
    MarketCap             float64 `json:"market_cap"`
    EpsTTM                float64 `json:"eps_ttm"`
    PeTTM                 float64 `json:"pe_ttm"`
    Pb                    float64 `json:"pb"`
    DividendYield         float64 `json:"dividend_yield"`
    Week52High            float64 `json:"week_52_high"`
    Week52Low             float64 `json:"week_52_low"`
    RevenueGrowthTTMYoy   float64 `json:"revenue_growth_ttm_yoy"`
    NetProfitMarginTTM    float64 `json:"net_profit_margin_ttm"`
    Beta                  float64 `json:"beta"`
    CurrentPrice           float64 `json:"current_price"`
    Score                  float64 `json:"score"`
}

// Custom unmarshal to handle "$" and convert to float64
func (s *StockInfo) UnmarshalJSON(data []byte) error {
    type Alias StockInfo
    aux := &struct {
        TargetFrom string `json:"target_from"`
        TargetTo   string `json:"target_to"`
        *Alias
    }{
        Alias: (*Alias)(s),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    s.TargetFrom, _ = parseDollarString(aux.TargetFrom)
    s.TargetTo, _ = parseDollarString(aux.TargetTo)
    return nil
}

func parseDollarString(val string) (float64, error) {
    val = strings.ReplaceAll(val, "$", "")
    val = strings.ReplaceAll(val, ",", "")
    return strconv.ParseFloat(val, 64)
}

func enrichWithFinnhub(stock *StockInfo) error {
    apiKey := os.Getenv("FINNHUB_API_KEY")
    cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
    client := finnhub.NewAPIClient(cfg).DefaultApi

    quote, _, err := client.Quote(context.Background()).Symbol(stock.Ticker).Execute()
    if err != nil {
        return err
    }
    if quote.C != nil {
        stock.CurrentPrice = float64(*quote.C) // Current price
    }

    basicFinancials, _, err := client.CompanyBasicFinancials(context.Background()).Symbol(stock.Ticker).Metric("all").Execute()
    if err == nil && basicFinancials.Metric != nil {
        m := *basicFinancials.Metric
        if v, ok := m["marketCapitalization"].(float64); ok {
            stock.MarketCap = v
        }
        if v, ok := m["epsTTM"].(float64); ok {
            stock.EpsTTM = v
        }
        if v, ok := m["peBasicTrailing12Months"].(float64); ok {
            stock.PeTTM = v
        }
        if v, ok := m["pbAnnual"].(float64); ok {
            stock.Pb = v
        }
        if v, ok := m["dividendYieldIndicatedAnnual"].(float64); ok {
            stock.DividendYield = v
        }
        if v, ok := m["52WeekHigh"].(float64); ok {
            stock.Week52High = v
        }
        if v, ok := m["52WeekLow"].(float64); ok {
            stock.Week52Low = v
        }
        if v, ok := m["revenueGrowthTTMYoy"].(float64); ok {
            stock.RevenueGrowthTTMYoy = v
        }
        if v, ok := m["netProfitMarginTTM"].(float64); ok {
            stock.NetProfitMarginTTM = v
        }
        if v, ok := m["beta"].(float64); ok {
            stock.Beta = v
        }
    }
    return nil
}

func scoreStock(stock StockInfo) float64 {
    score := 0.0

    // 1. Potencial de ganancia (target - actual) / actual
    if stock.CurrentPrice > 0 && stock.TargetTo > 0 {
        potential := (stock.TargetTo - stock.CurrentPrice) / stock.CurrentPrice
        score += potential * 100 // convertir a %
    }

    // 2. Precio objetivo anterior (TargetFrom): bonifica si el nuevo target es mayor
    if stock.TargetFrom > 0 && stock.TargetTo > stock.TargetFrom {
        score += (stock.TargetTo - stock.TargetFrom) * 0.1
    }

    // 3. Market Cap: bonifica empresas grandes, penaliza muy pequeñas
    if stock.MarketCap > 1e10 {
        score += 2
    } else if stock.MarketCap < 1e9 {
        score -= 2
    }

    // 4. EPS TTM: bonifica EPS positivo
    if stock.EpsTTM > 0 {
        score += stock.EpsTTM * 0.1
    }

    // 5. P/E TTM: menor a 20 es razonable
    if stock.PeTTM > 0 && stock.PeTTM < 20 {
        score += (20 - stock.PeTTM) * 0.5
    } else if stock.PeTTM <= 0 {
        score -= 5 // castigo por no tener ganancias
    }

    // 6. P/B: menor a 3 es razonable
    if stock.Pb > 0 && stock.Pb < 3 {
        score += (3 - stock.Pb) * 2
    }

    // 7. Dividend Yield: bonifica yield alto
    if stock.DividendYield > 1 {
        score += stock.DividendYield * 0.5
    }

    // 8. 52 Week High/Low: bonifica si el precio actual está más cerca del mínimo anual
    if stock.Week52High > 0 && stock.Week52Low > 0 && stock.CurrentPrice > 0 {
        rel := (stock.CurrentPrice - stock.Week52Low) / (stock.Week52High - stock.Week52Low)
        score += (1 - rel) * 2 // más cerca del mínimo, mejor
    }

    // 9. Revenue Growth YoY: bonifica crecimiento
    score += stock.RevenueGrowthTTMYoy * 0.2

    // 10. Net Profit Margin TTM: bonifica margen positivo
    if stock.NetProfitMarginTTM > 0 {
        score += stock.NetProfitMarginTTM * 0.2
    } else {
        score += stock.NetProfitMarginTTM * 0.1 // penaliza menos los negativos
    }

    // 11. Beta: ideal entre 0.8 y 1.2
    if stock.Beta >= 0.8 && stock.Beta <= 1.2 {
        score += 3
    } else if stock.Beta < 0.8 {
        score += 1
    } else {
        score -= 2
    }

    return score
}

func getStocksHandler(w http.ResponseWriter, r *http.Request) {

    // Allow CORS for local development
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight requests
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    conn, err := pgx.Connect(context.Background(), os.Getenv("COCKROACHDB_URL"))
    if err != nil {
        http.Error(w, "DB connection error", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    rows, err := conn.Query(context.Background(), `SELECT ticker, company, brokerage, stock_action, rating_from, rating_to, target_from, target_to, 
    stock_time, market_cap, eps_ttm, pe_ttm, pb, dividend_yield, week_52_high, week_52_low, revenue_growth_ttm_yoy, 
    net_profit_margin_ttm, beta , current_price, score FROM stock_info`)
    if err != nil {
        http.Error(w, "DB query error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var stocks []StockInfo
    for rows.Next() {
        var s StockInfo
        var stockTime time.Time
        err := rows.Scan(
            &s.Ticker,
            &s.Company,
            &s.Brokerage,
            &s.StockAction,
            &s.RatingFrom,
            &s.RatingTo,
            &s.TargetFrom,
            &s.TargetTo,
            &stockTime,
            &s.MarketCap,
            &s.EpsTTM,
            &s.PeTTM,
            &s.Pb,
            &s.DividendYield,
            &s.Week52High,
            &s.Week52Low,
            &s.RevenueGrowthTTMYoy,
            &s.NetProfitMarginTTM,
            &s.Beta,
            &s.CurrentPrice,
            &s.Score, // Asegúrate de que este campo esté en la consulta SQL
        )
        if err != nil {
            http.Error(w, "DB scan error", http.StatusInternalServerError)
            return
        }
        s.StockTime = stockTime.Format(time.RFC3339) // Convert to string
        stocks = append(stocks, s)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stocks)
}

func main() {

        // Cargar variables de entorno desde el archivo .env
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file") // Verificamos que las variables de entorno se hayan cargado correctamente.
        os.Exit(1)
    }

    apiURL := os.Getenv("API_URL") 
    apiKey := os.Getenv("API_KEY")
    
	// Creamos un cliente HTTP que se encargará de hacer la solicitud.
	client := &http.Client{}

	// Construimos una solicitud HTTP GET hacia la URL de la API.
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		// Si ocurre un error creando la solicitud, lo mostramos y detenemos el programa.
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// Añadimos el token de autorización en el header como tipo Bearer.
	req.Header.Add("Authorization", "Bearer "+apiKey)

	// Aunque no es obligatorio para GET, se indica que esperamos datos en formato JSON.
	req.Header.Add("Content-Type", "application/json")

	// Ejecutamos la solicitud usando el cliente HTTP.
	resp, err := client.Do(req)
	if err != nil {
		// Si hay un problema de red o la API no responde, mostramos el error y salimos.
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}

	// Nos aseguramos de cerrar el cuerpo de la respuesta al finalizar, para liberar recursos del sistema.
	defer resp.Body.Close()

    // Verificamos si la API respondió con éxito (200 OK)
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("API error: %s (status code: %d)\n", resp.Status, resp.StatusCode)
        os.Exit(1)
    }
    fmt.Println("¡Solicitud exitosa!")

	// Leemos todo el contenido del cuerpo de la respuesta y lo almacenamos como un array de bytes.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Si ocurre un error al leer el body (por ejemplo, por respuesta malformada), mostramos y salimos.
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	// Creamos una variable del tipo esperado para almacenar los datos del JSON.
	var apiResp ApiResponse

	// Parseamos el JSON (en formato []byte) y lo convertimos a nuestra estructura ApiResponse.
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// Si el JSON no se puede interpretar correctamente, lo mostramos para debug y salimos.
		fmt.Println("Error parsing JSON:", err)
		fmt.Println(string(body)) // Imprime la respuesta sin procesar para inspección manual.
		os.Exit(1)
	}

    for i := range apiResp.Items {
        err := enrichWithFinnhub(&apiResp.Items[i])
        if err != nil {
            fmt.Printf("Finnhub error for %s: %v\n", apiResp.Items[i].Ticker, err)
        }
    }
    // After enrichment:
    for i := range apiResp.Items {
        apiResp.Items[i].Score = scoreStock(apiResp.Items[i])
    }
	// Recorremos cada acción recibida del API y la imprimimos en un formato legible.
	for _, stock := range apiResp.Items {
		// Mostramos los detalles clave de cada recomendación de acción en una sola línea formateada.
		fmt.Printf(
			"Ticker: %s | Company: %s | Brokerage: %s | Action: %s | Rating: %s → %s | Target: %s → %s | Time: %s\n | Market Cap: %s | EPS: %s | P/E: %s | P/B: %s | Dividend Yield: %s | 52W High: %s | 52W Low: %s | Revenue Growth: %s | Net Profit Margin: %s | Beta: %s | Current Price: %s | Score: %.2f\n",
			stock.Ticker,
			stock.Company,
			stock.Brokerage,
			stock.StockAction,
			stock.RatingFrom,
			stock.RatingTo,
			stock.TargetFrom,
			stock.TargetTo,
			stock.StockTime,
            stock.MarketCap,
            stock.EpsTTM,
            stock.PeTTM,
            stock.Pb,
            stock.DividendYield,
            stock.Week52High,
            stock.Week52Low,
            stock.RevenueGrowthTTMYoy,
            stock.NetProfitMarginTTM,
            stock.Beta,
            stock.CurrentPrice,
            stock.Score,
		)
	}

    // --- CockroachDB connection setup ---
    // Establece la conexión a la base de datos CockroachDB usando el driver pgx.
    // La URL se obtiene desde una variable de entorno (COCKROACHDB_URL).
    conn, err := pgx.Connect(context.Background(), os.Getenv("COCKROACHDB_URL"))
    if err != nil {
        // Si ocurre un error al conectarse, lo mostramos y salimos del programa.
        fmt.Println("Error connecting to CockroachDB:", err)
        os.Exit(1)
    }
    fmt.Println("Conexión a CockroachDB exitosa.")
    // Defer garantiza que la conexión se cerrará al final de main, liberando recursos.
    defer conn.Close(context.Background())

    // Ejecuta una sentencia SQL para crear una tabla llamada 'stock_info' si no existe aún.
    _, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS stock_info (
            ticker STRING,
            stock_time TIMESTAMPTZ,
            company STRING,
            brokerage STRING,
            stock_action STRING,
            rating_from STRING,
            rating_to STRING,
            target_from FLOAT8,
            target_to FLOAT8,
            market_cap FLOAT8,
            eps_ttm FLOAT8,
            pe_ttm FLOAT8,
            pb FLOAT8,
            dividend_yield FLOAT8,
            week_52_high FLOAT8,
            week_52_low FLOAT8,
            revenue_growth_ttm_yoy FLOAT8,
            net_profit_margin_ttm FLOAT8,
            beta FLOAT8,
            current_price FLOAT8,
            score FLOAT8,
            PRIMARY KEY (ticker, stock_time)
            );
    `)
    if err != nil {
        // Si hubo un error creando la tabla, se muestra y el programa termina
        fmt.Println("Error creating table:", err)
        os.Exit(1)
    }

    // Recorremos cada acción bursátil recibida del API
    for _, stock := range apiResp.Items {
        // Ejecutamos una sentencia UPSERT para guardar cada acción en la tabla 'stock_info'
        _, err := conn.Exec(context.Background(),
            `UPSERT INTO stock_info (
                ticker, company, brokerage, stock_action,
                rating_from, rating_to, target_from, target_to, stock_time,
                market_cap, eps_ttm, pe_ttm, pb, dividend_yield,
                week_52_high, week_52_low, revenue_growth_ttm_yoy,
                net_profit_margin_ttm, beta, current_price, score
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`, // Usamos placeholders para evitar SQL injection
            stock.Ticker,       // $1 → símbolo de la acción
            stock.Company,      // $2 → nombre de la empresa
            stock.Brokerage,    // $3 → casa de bolsa
            stock.StockAction,       // $4 → acción recomendada (ej: Buy)
            stock.RatingFrom,   // $5 → calificación anterior
            stock.RatingTo,     // $6 → nueva calificación
            stock.TargetFrom,   // $7 → precio objetivo anterior
            stock.TargetTo,     // $8 → nuevo precio objetivo
            stock.StockTime,         // $9 → marca de tiempo
            stock.MarketCap,    // $10 → capitalización de mercado  
            stock.EpsTTM,       // $11 → ganancias por acción (EPS) en los últimos 12 meses
            stock.PeTTM,        // $12 → relación precio/ganancias (P/E) en los últimos 12 meses
            stock.Pb,           // $13 → relación precio/valor contable (P/B        
            stock.DividendYield, // $14 → rendimiento por dividendo
            stock.Week52High,   // $15 → máximo de 52 semanas
            stock.Week52Low,    // $16 → mínimo de 52 semanas
            stock.RevenueGrowthTTMYoy, // $17 → crecimiento de ingresos interanual
            stock.NetProfitMarginTTM, // $18 → margen de beneficio neto
            stock.Beta,        // $19 → beta (volatilidad relativa)
            stock.CurrentPrice, // $20 → precio actual
            stock.Score,         // $21 → puntaje calculado
        )

        // Si ocurre un error al insertar esta fila, lo mostramos
        if err != nil {
            fmt.Printf("Error inserting stock %s: %v\n", stock.Ticker, err)
        }
    }

    // Mensaje final confirmando que los datos fueron insertados sin errores críticos
    fmt.Println("Datos almacenados en CockroachDB correctamente.")

	http.HandleFunc("/stocks", getStocksHandler)
    fmt.Println("Backend API running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
