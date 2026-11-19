package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"prj_withdrawal/headers"
	"sync"
)

// ==========================================
// ========== INICIO ACÁPITE 4 ==============
// ==========================================

var (
	cuentaGlobal *headers.Cuenta
	mu           sync.Mutex
)

type OperacionRequest struct {
	Monto float64 `json:"monto"`
}

type ReiniciarRequest struct {
	Titular      string  `json:"titular"`
	SaldoInicial float64 `json:"saldo_inicial"`
}

type EstadoResponse struct {
	Titular   string                `json:"titular"`
	Saldo     float64               `json:"saldo"`
	Estado    string                `json:"estado"`
	Mensaje   string                `json:"mensaje"`
	Historial []headers.Transaccion `json:"historial"`
}

func responderJSON(w http.ResponseWriter, codigo int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(data)
}

func handleEstado(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	responderJSON(w, http.StatusOK, EstadoResponse{
		Titular:   cuentaGlobal.Nombre(),
		Saldo:     cuentaGlobal.ObtenerSaldo(),
		Estado:    "OK",
		Mensaje:   "Cuenta cargada correctamente",
		Historial: cuentaGlobal.ObtenerHistorial(),
	})
}

func handleRetirar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var req OperacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderJSON(w, http.StatusBadRequest, EstadoResponse{Estado: "ERROR", Mensaje: "Monto inválido"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	exito, err := cuentaGlobal.Retirar(req.Monto)
	cuentaGlobal.ExportarLogJSON("transacciones.json")

	estado := "EXITO"
	mensaje := fmt.Sprintf("Retiro de S/ %.2f procesado correctamente.", req.Monto)
	if err != nil {
		estado = "ERROR"
		mensaje = err.Error()
	} else if !exito {
		estado = "FALLO"
		mensaje = "Saldo insuficiente. Intente con un monto menor."
	}

	responderJSON(w, http.StatusOK, EstadoResponse{
		Titular:   cuentaGlobal.Nombre(),
		Saldo:     cuentaGlobal.ObtenerSaldo(),
		Estado:    estado,
		Mensaje:   mensaje,
		Historial: cuentaGlobal.ObtenerHistorial(),
	})
}

func handleAbonar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var req OperacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responderJSON(w, http.StatusBadRequest, EstadoResponse{Estado: "ERROR", Mensaje: "Monto inválido"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, err := cuentaGlobal.Abonar(req.Monto)
	cuentaGlobal.ExportarLogJSON("transacciones.json")

	estado := "EXITO"
	mensaje := fmt.Sprintf("Abono de S/ %.2f procesado correctamente.", req.Monto)
	if err != nil {
		estado = "ERROR"
		mensaje = err.Error()
	}

	responderJSON(w, http.StatusOK, EstadoResponse{
		Titular:   cuentaGlobal.Nombre(),
		Saldo:     cuentaGlobal.ObtenerSaldo(),
		Estado:    estado,
		Mensaje:   mensaje,
		Historial: cuentaGlobal.ObtenerHistorial(),
	})
}

func handleReiniciar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var req ReiniciarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Titular == "" {
		req.Titular = "Marco Gomez"
		req.SaldoInicial = 200.00
	}

	mu.Lock()
	defer mu.Unlock()

	cuentaGlobal = headers.NewCuenta(req.Titular, req.SaldoInicial)
	cuentaGlobal.ExportarLogJSON("transacciones.json")

	responderJSON(w, http.StatusOK, EstadoResponse{
		Titular:   cuentaGlobal.Nombre(),
		Saldo:     cuentaGlobal.ObtenerSaldo(),
		Estado:    "OK",
		Mensaje:   "Sesión reiniciada correctamente.",
		Historial: cuentaGlobal.ObtenerHistorial(),
	})
}

// ==========================================
// =========== FIN ACÁPITE 4 ================
// ==========================================

func main() {
	ejecutarBDD := flag.Bool("bdd", false, "Ejecuta los escenarios BDD/TDD por consola")
	puerto := flag.String("port", "8080", "Puerto del servidor web")
	flag.Parse()

	// Modo CLI / Pruebas BDD heredado de los Acápites 1, 2 y 3
	if *ejecutarBDD {
		fmt.Println("=== EJECUTANDO PRUEBAS BDD (Gherkin v1.0.0) ===")
		headers.TestEscenario1Exito("Alice", 234.56, 200.00)
		headers.TestEscenario2Fallo("Hamza", 198.76, 200.00)
		headers.TestEscenario4Abono("Carlos", 100.00, 50.00)
		headers.TestEscenario5IntegridadJSON("Marco", 200.00)
		os.Exit(0)
	}

	// Inicialización de la cuenta para la App Web
	cuentaGlobal = headers.NewCuenta("Marco Gomez", 200.00)
	cuentaGlobal.ExportarLogJSON("transacciones.json")

	// Rutas de la API REST
	http.HandleFunc("/api/estado", handleEstado)
	http.HandleFunc("/api/retirar", handleRetirar)
	http.HandleFunc("/api/abonar", handleAbonar)
	http.HandleFunc("/api/reiniciar", handleReiniciar)

	// Servir archivos estáticos del Frontend (Capa de Presentación)
	http.Handle("/", http.FileServer(http.Dir("./public")))

	fmt.Printf("Servidor ATM Web activo en http://localhost:%s\n", *puerto)
	fmt.Println("Presiona Ctrl + C para detener el servidor.")
	if err := http.ListenAndServe(":"+*puerto, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
		os.Exit(1)
	}
}