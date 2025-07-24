// Declaramos que este archivo pertenece al paquete principal (main), obligatorio para ejecutar un programa Go.
package main

// Importamos los paquetes necesarios para hacer solicitudes HTTP, leer datos, manejar errores, parsear JSON, etc.
import (
	"encoding/json" // Permite convertir entre JSON y estructuras Go.
	"fmt"           // Sirve para imprimir texto en consola.
	"io"            // Proporciona utilidades para leer el cuerpo de la respuesta HTTP.
	"net/http"      // Permite crear clientes HTTP y hacer solicitudes.
	"os"            // Permite interactuar con el sistema operativo (ej: salir del programa con os.Exit).

	"github.com/joho/godotenv" // opcional, para cargar .env en local
	// This will create a go.mod file and allow you to use external packages like github.com/joho/godotenv.
	//go mod init stockradar
	//go get github.com/joho/godotenv
    "context"
    "github.com/jackc/pgx/v5" // Importa el paquete pgx para manejar conexiones a bases de datos PostgreSQL.
    //go get github.com/jackc/pgx/v5

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
	Action      string `json:"action"`       // Tipo de recomendación (ej: "Buy", "Hold", etc).
	RatingFrom  string `json:"rating_from"`  // Calificación anterior del bróker.
	RatingTo    string `json:"rating_to"`    // Nueva calificación otorgada.
	TargetFrom  string `json:"target_from"`  // Precio objetivo anterior.
	TargetTo    string `json:"target_to"`    // Precio objetivo nuevo.
	Time        string `json:"time"`         // Momento en que se emitió la recomendación.
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

	// Recorremos cada acción recibida del API y la imprimimos en un formato legible.
	for _, stock := range apiResp.Items {
		// Mostramos los detalles clave de cada recomendación de acción en una sola línea formateada.
		fmt.Printf(
			"Ticker: %s | Company: %s | Brokerage: %s | Action: %s | Rating: %s → %s | Target: %s → %s | Time: %s\n",
			stock.Ticker,
			stock.Company,
			stock.Brokerage,
			stock.Action,
			stock.RatingFrom,
			stock.RatingTo,
			stock.TargetFrom,
			stock.TargetTo,
			stock.Time,
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
            ticker STRING,        -- Símbolo bursátil (ej: AAPL, TSLA)
            company STRING,       -- Nombre de la empresa
            brokerage STRING,     -- Casa de bolsa que emitió la recomendación
            action STRING,        -- Acción recomendada (ej: Buy, Hold, Sell)
            rating_from STRING,   -- Calificación anterior (ej: Neutral)
            rating_to STRING,     -- Nueva calificación (ej: Buy)
            target_from STRING,   -- Precio objetivo anterior
            target_to STRING,     -- Precio objetivo actualizado
            time STRING           -- Marca de tiempo (puede ser fecha o fecha-hora)
        )
    `)
    if err != nil {
        // Si hubo un error creando la tabla, se muestra y el programa termina
        fmt.Println("Error creating table:", err)
        os.Exit(1)
    }

    // Recorremos cada acción bursátil recibida del API
    for _, stock := range apiResp.Items {
        // Ejecutamos una sentencia INSERT para guardar cada acción en la tabla 'stock_info'
        _, err := conn.Exec(context.Background(),
            `INSERT INTO stock_info (
                ticker, company, brokerage, action,
                rating_from, rating_to, target_from, target_to, time
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, // Usamos placeholders para evitar SQL injection
            stock.Ticker,       // $1 → símbolo de la acción
            stock.Company,      // $2 → nombre de la empresa
            stock.Brokerage,    // $3 → casa de bolsa
            stock.Action,       // $4 → acción recomendada (ej: Buy)
            stock.RatingFrom,   // $5 → calificación anterior
            stock.RatingTo,     // $6 → nueva calificación
            stock.TargetFrom,   // $7 → precio objetivo anterior
            stock.TargetTo,     // $8 → nuevo precio objetivo
            stock.Time,         // $9 → marca de tiempo
        )

        // Si ocurre un error al insertar esta fila, lo mostramos
        if err != nil {
            fmt.Printf("Error inserting stock %s: %v\n", stock.Ticker, err)
        }
    }

    // Mensaje final confirmando que los datos fueron insertados sin errores críticos
    fmt.Println("Datos almacenados en CockroachDB correctamente.")
}
