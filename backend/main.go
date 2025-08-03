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

	"context" // Permite manejar el contexto de ejecución, útil para operaciones con base de datos y APIs externas.

	"github.com/joho/godotenv" // Para cargar variables de entorno desde un archivo .env en desarrollo.

	"time" // Permite trabajar con fechas y horas.

	"github.com/jackc/pgx/v5" // Driver para conectar y operar con bases de datos PostgreSQL/CockroachDB.

	"math"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2" // SDK para consumir la API de Finnhub.
)

// Estructura que representa la respuesta principal de la API externa.
// Contiene un arreglo de acciones (StockInfo) y un campo opcional para paginación.
type ApiResponse struct {
    Items    []StockInfo `json:"items"`      // Lista de acciones recibidas del API.
    NextPage string      `json:"next_page"`  // Paginación, por si hay más datos (opcional).
}

// Estructura que representa cada acción (stock) con todos los campos financieros relevantes.
// Los tags json permiten mapear los campos del JSON con los campos de la estructura Go.
type StockInfo struct {
    Ticker      string  `json:"ticker"`        // Símbolo bursátil de la acción.
    Company     string  `json:"company"`       // Nombre de la empresa.
    Brokerage   string  `json:"brokerage"`     // Nombre del bróker que reporta la recomendación.
    StockAction string  `json:"action"`        // Tipo de recomendación (ej: "Buy", "Hold", etc).
    RatingFrom  string  `json:"rating_from"`   // Calificación anterior del bróker.
    RatingTo    string  `json:"rating_to"`     // Nueva calificación otorgada.
    TargetFrom  float64 `json:"target_from"`   // Precio objetivo anterior.
    TargetTo    float64 `json:"target_to"`     // Precio objetivo nuevo.
    StockTime   string  `json:"time"`          // Momento en que se emitió la recomendación.
    MarketCap             float64 `json:"market_cap"`              // Capitalización de mercado.
    EpsTTM                float64 `json:"eps_ttm"`                 // Ganancias por acción (últimos 12 meses).
    PeTTM                 float64 `json:"pe_ttm"`                  // Relación precio/ganancias (últimos 12 meses).
    Pb                    float64 `json:"pb"`                      // Relación precio/valor contable.
    DividendYield         float64 `json:"dividend_yield"`          // Rendimiento por dividendo.
    Week52High            float64 `json:"week_52_high"`            // Máximo de 52 semanas.
    Week52Low             float64 `json:"week_52_low"`             // Mínimo de 52 semanas.
    RevenueGrowthTTMYoy   float64 `json:"revenue_growth_ttm_yoy"`  // Crecimiento de ingresos interanual.
    NetProfitMarginTTM    float64 `json:"net_profit_margin_ttm"`   // Margen de beneficio neto.
    Beta                  float64 `json:"beta"`                    // Beta (volatilidad relativa).
    CurrentPrice          float64 `json:"current_price"`           // Precio actual de la acción.
    Score                 float64 `json:"score"`                   // Puntaje calculado según criterios cuantitativos.
    Normalized            float64 `json:"normalized"` 
}


// UnmarshalJSON personalizado para convertir los campos target_from y target_to de string (con $) a float64.
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

// Función auxiliar para eliminar el símbolo $ y convertir a float64.
func parseDollarString(val string) (float64, error) {
    val = strings.ReplaceAll(val, "$", "")
    val = strings.ReplaceAll(val, ",", "")
    return strconv.ParseFloat(val, 64)
}

// Enriquecimiento de datos usando la API de Finnhub.
// Obtiene precio actual y métricas financieras clave para cada acción.
func enrichWithFinnhub(stock *StockInfo) error {
    apiKey := os.Getenv("FINNHUB_API_KEY")
    cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
    client := finnhub.NewAPIClient(cfg).DefaultApi

    // Consulta el precio actual de la acción.
    quote, _, err := client.Quote(context.Background()).Symbol(stock.Ticker).Execute()
    if err != nil {
        return err
    }
    if quote.C != nil {
        stock.CurrentPrice = float64(*quote.C) // Precio actual
    }

    // Consulta los datos financieros básicos.
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

func NormalizeScores(stocks []*StockInfo) {
    // 1. Transformar valores negativos
    for _, s := range stocks {
        if s.Score < 0 {
            abs := math.Abs(s.Score)
            log10 := math.Log10(abs)
            s.Score = 1 / math.Pow(10, log10)
        }
    }

    // 2. Encontrar min y max (después de transformar negativos)
    minScore := math.MaxFloat64
    maxScore := -math.MaxFloat64
    for _, s := range stocks {
        if s.Score < minScore {
            minScore = s.Score
        }
        if s.Score > maxScore {
            maxScore = s.Score
        }
    }

    // 3. Normalizar y actualizar el campo Normalized
    for _, s := range stocks {
        if maxScore != minScore {
            norm := (s.Score - minScore) / (maxScore - minScore)
            if norm < 0 {
                s.Normalized = 0
            } else {
                s.Normalized = norm
            }
        } else {
            s.Normalized = 0.5
        }
    }
}

// Calcula un puntaje cuantitativo para cada acción usando todos los campos numéricos relevantes.
func scoreStock(stock StockInfo) float64 {
    score := 0.0

    // 1. Potencial de ganancia (target - actual) / actual
    if stock.CurrentPrice > 0 && stock.TargetTo > 0 {
        potential := (stock.TargetTo - stock.CurrentPrice) / stock.CurrentPrice
        score += potential * 150 // convertir a %
    }

    // 2. Precio objetivo anterior (TargetFrom): bonifica si el nuevo target es mayor
    if stock.TargetFrom > 0 && stock.TargetTo > stock.TargetFrom {
        score += (stock.TargetTo - stock.TargetFrom) * 0.2
    }

    // 3. Market Cap: bonifica empresas grandes, penaliza muy pequeñas
    if stock.MarketCap > 1e10 {
        score += 5
    } else if stock.MarketCap < 1e9 {
        score -= 5
    }

    // 4. EPS TTM: bonificación más fuerte si es positivo
    if stock.EpsTTM > 0 {
        score += stock.EpsTTM * 0.5
    } else {
        score += stock.EpsTTM * 1.5 // castigo mayor si es negativo
    }

    // 5. P/E TTM: menor a 20 es razonable
    if stock.PeTTM > 0 && stock.PeTTM < 20 {
        score += (20 - stock.PeTTM) * 1.0
    } else if stock.PeTTM <= 0 {
        score -= 10 // castigo por no tener ganancias
    }

    // 6. P/B: mayor penalización a valores altos
    if stock.Pb > 0 && stock.Pb < 3 {
        score += (3 - stock.Pb) * 3
    } else if stock.Pb >= 3 {
        score -= (stock.Pb - 3) * 1.5
    }

    // 7. Dividend Yield: bonifica yield alto
    if stock.DividendYield > 1 {
        score += stock.DividendYield * 1.0
    }

    // 8. 52 Week High/Low: bonifica si el precio actual está más cerca del mínimo anual
    if stock.Week52High > 0 && stock.Week52Low > 0 && stock.CurrentPrice > 0 {
        rel := (stock.CurrentPrice - stock.Week52Low) / (stock.Week52High - stock.Week52Low)
        score += (1 - rel) * 5 // más cerca del mínimo, mejor
    }

    // 9. Revenue Growth YoY: bonifica crecimiento
    score += stock.RevenueGrowthTTMYoy * 0.6

    // 10. Net Profit Margin TTM: bonifica margen positivo
    if stock.NetProfitMarginTTM > 0 {
        score += stock.NetProfitMarginTTM * 0.4
    } else {
        score += stock.NetProfitMarginTTM * 0.2 // penaliza menos los negativos
    }

    // 11. Beta: ideal entre 0.8 y 1.2
    if stock.Beta >= 0.8 && stock.Beta <= 1.2 {
        score += 5
    } else if stock.Beta < 0.8 {
        score += 2
    } else {
        score -= (stock.Beta - 1.2) * 3 // castigo proporcional
    }

    return score
}

// Handler HTTP para exponer los datos almacenados en la base de datos vía API REST.
func getStocksHandler(w http.ResponseWriter, r *http.Request) {
    // Permitir CORS para desarrollo local
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Manejar preflight requests
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    // Conexión a la base de datos
    conn, err := pgx.Connect(context.Background(), os.Getenv("COCKROACHDB_URL"))
    if err != nil {
        http.Error(w, "DB connection error", http.StatusInternalServerError)
        return
    }
    defer conn.Close(context.Background())

    // Consulta SQL para obtener todos los campos de la tabla stock_info
    rows, err := conn.Query(context.Background(), `SELECT ticker, company, brokerage, stock_action, rating_from, rating_to, target_from, target_to, 
        stock_time, market_cap, eps_ttm, pe_ttm, pb, dividend_yield, week_52_high, week_52_low, revenue_growth_ttm_yoy, 
        net_profit_margin_ttm, beta , current_price, score, normalized FROM stock_info`)
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
            &s.Score,
            &s.Normalized,
        )
        if err != nil {
            http.Error(w, "DB scan error", http.StatusInternalServerError)
            return
        }
        s.StockTime = stockTime.Format(time.RFC3339) // Convertir a string ISO
        stocks = append(stocks, s)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stocks)
}

// Función principal del backend.
// Orquesta la carga de datos, enriquecimiento, almacenamiento y exposición vía API.
func main() {
    // Cargar variables de entorno desde el archivo .env
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
        os.Exit(1)
    }

    // Leer la URL y el API key de la API externa desde variables de entorno
    apiURL := os.Getenv("API_URL")
    apiKey := os.Getenv("API_KEY")

    // Crear cliente HTTP para consumir la API externa
	client := &http.Client{}

    // Construir la solicitud HTTP GET
    req, err := http.NewRequest("GET", apiURL, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        os.Exit(1)
    }

    // Añadir el token de autorización y el tipo de contenido
    req.Header.Add("Authorization", "Bearer "+apiKey)
    req.Header.Add("Content-Type", "application/json")

    // Ejecutar la solicitud HTTP
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error making request:", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    // Verificar que la respuesta sea exitosa
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("API error: %s (status code: %d)\n", resp.Status, resp.StatusCode)
        os.Exit(1)
    }
    fmt.Println("¡Solicitud exitosa!")

    // Leer el cuerpo de la respuesta
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response:", err)
        os.Exit(1)
    }

    // Parsear el JSON recibido a la estructura ApiResponse
    var apiResp ApiResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        fmt.Println("Error parsing JSON:", err)
        fmt.Println(string(body))
        os.Exit(1)
    }

    // Enriquecer cada acción con datos de Finnhub
    for i := range apiResp.Items {
        err := enrichWithFinnhub(&apiResp.Items[i])
        if err != nil {
            fmt.Printf("Finnhub error for %s: %v\n", apiResp.Items[i].Ticker, err)
        }
    }

    // Calcular el puntaje para cada acción
    for i := range apiResp.Items {
        apiResp.Items[i].Score = scoreStock(apiResp.Items[i])
    }

    for _, s := range apiResp.Items {
    fmt.Println("Score antes de normalizar:", s.Score)
    }
    // Prepara un slice de punteros
    var stockPtrs []*StockInfo
    for i := range apiResp.Items {
        stockPtrs = append(stockPtrs, &apiResp.Items[i])
    }
    NormalizeScores(stockPtrs)

    // Imprimir en consola los datos enriquecidos y el puntaje
    for _, stock := range apiResp.Items {
        fmt.Printf(
            "Ticker: %s | Company: %s | Brokerage: %s | Action: %s | Rating: %s → %s | Target: %.2f → %.2f | Time: %s\n | Market Cap: %.2f | EPS: %.2f | P/E: %.2f | P/B: %.2f | Dividend Yield: %.2f | 52W High: %.2f | 52W Low: %.2f | Revenue Growth: %.2f | Net Profit Margin: %.2f | Beta: %.2f | Current Price: %.2f | Score: %.2f | Normalized: %.2f\n",
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
            stock.Normalized,
        )
    }

    // --- Conexión y almacenamiento en CockroachDB ---
    // Conectar a la base de datos CockroachDB usando la URL de entorno
    conn, err := pgx.Connect(context.Background(), os.Getenv("COCKROACHDB_URL"))
    if err != nil {
        fmt.Println("Error connecting to CockroachDB:", err)
        os.Exit(1)
    }
    fmt.Println("Conexión a CockroachDB exitosa.")
    defer conn.Close(context.Background())

    // Crear la tabla stock_info si no existe
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
            normalized FLOAT8,
            PRIMARY KEY (ticker, stock_time)
        );
    `)
    if err != nil {
        fmt.Println("Error creating table:", err)
        os.Exit(1)
    }

    // Insertar o actualizar cada acción en la base de datos
    for _, stock := range apiResp.Items {
        _, err := conn.Exec(context.Background(),
            `UPSERT INTO stock_info (
                ticker, company, brokerage, stock_action,
                rating_from, rating_to, target_from, target_to, stock_time,
                market_cap, eps_ttm, pe_ttm, pb, dividend_yield,
                week_52_high, week_52_low, revenue_growth_ttm_yoy,
                net_profit_margin_ttm, beta, current_price, score, normalized
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`,
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
            stock.Normalized,
        )
        if err != nil {
            fmt.Printf("Error inserting stock %s: %v\n", stock.Ticker, err)
        }
    }
    for _, stock := range apiResp.Items {
        fmt.Printf("Ticker: %s | Score: %.2f | Normalized: %.4f\n", stock.Ticker, stock.Score, stock.Normalized)
    }

    fmt.Println("Datos almacenados en CockroachDB correctamente.")

    // Iniciar el servidor HTTP para exponer la API REST en /stocks
    http.HandleFunc("/stocks", getStocksHandler)
    fmt.Println("Backend API running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
