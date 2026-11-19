package headers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// ==========================================
// ========== INICIO ACÁPITE 3 ==============
// ==========================================

// Transaccion representa el registro de auditoría de cada operación (US-04).
type Transaccion struct {
	Fecha      string  `json:"fecha"`
	Titular    string  `json:"titular"`
	Tipo       string  `json:"tipo"`       // RETIRO o ABONO
	Monto      float64 `json:"monto"`
	Estado     string  `json:"estado"`     // EXITO, FALLO o ERROR
	Detalle    string  `json:"detalle"`
	SaldoFinal float64 `json:"saldo_final"`
}

// ==========================================
// =========== FIN ACÁPITE 3 ================
// ==========================================

// Cuenta representa la entidad de negocio descrita en el Feature Gherkin.
type Cuenta struct {
	saldo     float64
	nombre    string
	historial []Transaccion // Acápite 3: historial para auditoría JSON
}

// NewCuenta es el constructor equivalente a Cuenta(nombre, saldoInicial).
func NewCuenta(nombre string, saldoInicial float64) *Cuenta {
	return &Cuenta{
		nombre:    nombre,
		saldo:     saldoInicial,
		historial: []Transaccion{},
	}
}

// ObtenerSaldo devuelve el saldo actual.
func (c *Cuenta) ObtenerSaldo() float64 {
	return c.saldo
}

// Nombre devuelve el nombre del titular (útil para logs).
func (c *Cuenta) Nombre() string {
	return c.nombre
}

// ==========================================
// ========== INICIO ACÁPITE 3 ==============
// ==========================================

// registrarTransaccion guarda cada intento (éxito, fallo o error) en el historial.
func (c *Cuenta) registrarTransaccion(tipo string, monto float64, estado string, detalle string) {
	tx := Transaccion{
		Fecha:      time.Now().Format("2006-01-02 15:04:05"),
		Titular:    c.nombre,
		Tipo:       tipo,
		Monto:      monto,
		Estado:     estado,
		Detalle:    detalle,
		SaldoFinal: c.saldo,
	}
	c.historial = append(c.historial, tx)
}

// ExportarLogJSON convierte el historial a formato JSON y lo guarda en un archivo local.
func (c *Cuenta) ExportarLogJSON(nombreArchivo string) (string, error) {
	datosJSON, err := json.MarshalIndent(c.historial, "", "  ")
	if err != nil {
		return "", err
	}
	err = os.WriteFile(nombreArchivo, datosJSON, 0644)
	if err != nil {
		return "", err
	}
	return string(datosJSON), nil
}

// ==========================================
// =========== FIN ACÁPITE 3 ================
// ==========================================

/*
Retirar intenta retirar una cantidad.
Regla: No se puede retirar más de lo que hay en la cuenta.
*/
func (c *Cuenta) Retirar(cantidad float64) (bool, error) {
	if cantidad < 0 {
		c.registrarTransaccion("RETIRO", cantidad, "ERROR", "No se pueden retirar cantidades negativas")
		return false, errors.New("no se pueden retirar cantidades negativas")
	}
	if cantidad > c.saldo {
		c.registrarTransaccion("RETIRO", cantidad, "FALLO", "Saldo insuficiente")
		return false, nil // Escenario: Withdrawal declined
	}
	c.saldo -= cantidad
	c.registrarTransaccion("RETIRO", cantidad, "EXITO", "Retiro procesado correctamente")
	return true, nil // Escenario: Successful withdrawal
}

// ==========================================
// ========== INICIO ACÁPITE 2 ==============
// ==========================================
/*
Abonar intenta depositar una cantidad en la cuenta.
Regla: No se pueden abonar cantidades negativas.
*/
func (c *Cuenta) Abonar(cantidad float64) (bool, error) {
	if cantidad < 0 {
		c.registrarTransaccion("ABONO", cantidad, "ERROR", "No se pueden abonar cantidades negativas")
		return false, errors.New("no se pueden abonar cantidades negativas")
	}

	c.saldo += cantidad
	c.registrarTransaccion("ABONO", cantidad, "EXITO", "Abono procesado correctamente")
	return true, nil // Escenario: Successful deposit
}

// ==========================================
// =========== FIN ACÁPITE 2 ================
// ==========================================

// String implementa fmt.Stringer para depuración
func (c *Cuenta) String() string {
	return fmt.Sprintf("Cuenta{%s, saldo=%.2f}", c.nombre, c.saldo)
}

// ==========================================
// ========== INICIO ACÁPITE 4 ==============
// ==========================================

// ObtenerHistorial devuelve la lista completa de transacciones registradas para la API Web.
func (c *Cuenta) ObtenerHistorial() []Transaccion {
	return c.historial
}

// ==========================================
// =========== FIN ACÁPITE 4 ================
// ==========================================	