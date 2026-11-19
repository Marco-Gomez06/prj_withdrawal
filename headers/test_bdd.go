package headers

import (
	"fmt"
	"math"
)

// imprimirSaldo imprime el saldo formateado a 2 decimales.
func imprimirSaldo(cuenta *Cuenta) {
	fmt.Printf("Saldo de %.2f \n", cuenta.ObtenerSaldo())
}

// ----------------------------------------------------------------------
// Validación de Escenarios Gherkin Dinámicos
// ----------------------------------------------------------------------

// TestEscenario1Exito: Simula un retiro exitoso con datos dinámicos
func TestEscenario1Exito(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 1: Successful withdrawal within balance (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if exito {
		fmt.Println("Resultado: ÉXITO. La retirada fue procesada.")
		imprimirSaldo(cuenta)
		saldoEsperado := saldoInicial - cantidadRetirar
		if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
			fmt.Printf("Validación: El saldo restante es correcto (%.2f).\n\n", saldoEsperado)
		} else {
			fmt.Println("ERROR: El saldo restante no coincide.")
			fmt.Println()
		}
	} else {
		fmt.Println("Resultado: FALLO. La retirada fue rechazada (incorrecto).")
		fmt.Println()
	}
}

// TestEscenario2Fallo: Simula un retiro rechazado con datos dinámicos
func TestEscenario2Fallo(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 2: Declined withdrawal in excess of balance (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta retirar %.2f\n", nombre, cantidadRetirar)
	exito, err := cuenta.Retirar(cantidadRetirar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: FALLO. Error inesperado: %v\n\n", err)
		return
	}

	if !exito {
		fmt.Println("Resultado: ÉXITO. La retirada fue rechazada correctamente.")
		imprimirSaldo(cuenta)
		if math.Abs(cuenta.ObtenerSaldo()-saldoInicial) < 0.001 {
			fmt.Printf("Validación: El saldo no ha cambiado (%.2f).\n\n", saldoInicial)
		} else {
			fmt.Println("ERROR: El saldo se modificó incorrectamente.")
			fmt.Println()
		}
	} else {
		fmt.Println("Resultado: FALLO. La retirada fue procesada (debería haber sido rechazada).")
		fmt.Println()
	}
}

// TestEscenario3Error: Simula un error de validación (ej. cantidad negativa)
func TestEscenario3Error(nombre string, saldoInicial float64, cantidadRetirar float64) {
	fmt.Printf("--- Escenario 3: Validation Error (%s) ---\n", nombre)
	cuenta := NewCuenta(nombre, saldoInicial)

	exito, err := cuenta.Retirar(cantidadRetirar)
	if err != nil {
		fmt.Printf("Resultado: ÉXITO. El sistema atrapó el error correctamente: %v\n", err)
		imprimirSaldo(cuenta)
		fmt.Println()
	} else if !exito {
		fmt.Println("Resultado: RECHAZADO. Fondos insuficientes (no es un error de validación).")
		fmt.Println()
	} else {
		fmt.Println("Resultado: FALLO. Se procesó un retiro que debía dar error.")
		fmt.Println()
	}
}

// ==========================================
// ========== INICIO ACAPITE 2 ==============
// ==========================================

// TestEscenario4Abono: Simula un abono exitoso con datos dinámicos
func TestEscenario4Abono(nombre string, saldoInicial float64, cantidadAbonar float64) {
	fmt.Printf("--- Escenario 4: Successful deposit (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When
	fmt.Printf("When: %s intenta abonar %.2f\n", nombre, cantidadAbonar)
	exito, err := cuenta.Abonar(cantidadAbonar)

	// Then
	if err != nil {
		fmt.Printf("Resultado: FALLO. Error de validación: %v\n\n", err)
		return
	}

	if exito {
		fmt.Println("Resultado: ÉXITO. El abono fue procesado.")
		imprimirSaldo(cuenta)
		saldoEsperado := saldoInicial + cantidadAbonar
		if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
			fmt.Printf("Validación: El nuevo saldo es correcto (%.2f).\n\n", saldoEsperado)
		} else {
			fmt.Println("ERROR: El saldo restante no coincide.")
			fmt.Println()
		}
	} else {
		fmt.Println("Resultado: FALLO. El abono fue rechazado (incorrecto).")
		fmt.Println()
	}
}

// ==========================================
// =========== FIN ACAPITE 2 ================
// ==========================================

// ==========================================
// ========== INICIO ACÁPITE 3 ==============
// ==========================================

// TestEscenario5IntegridadJSON: Valida US-04 (atomicidad y registro JSON de éxito, fallo y error)
func TestEscenario5IntegridadJSON(nombre string, saldoInicial float64) {
	fmt.Printf("--- Escenario 5: Integridad de Transacción y Log JSON (%s) ---\n", nombre)

	// Given
	cuenta := NewCuenta(nombre, saldoInicial)
	fmt.Printf("Given: %s tiene %.2f en su cuenta\n", nombre, saldoInicial)

	// When: ejecutamos 3 operaciones (Éxito, Fallo y Error)
	fmt.Println("When: realiza un retiro válido (50.00), uno sin fondos (500.00) y uno negativo (-20.00)")
	cuenta.Retirar(50.00)  // EXITO -> saldo baja a 150.00
	cuenta.Retirar(500.00) // FALLO -> saldo sigue en 150.00
	cuenta.Retirar(-20.00) // ERROR -> saldo sigue en 150.00

	// Then: validamos atomicidad del saldo y generamos el archivo JSON
	saldoEsperado := saldoInicial - 50.00
	if math.Abs(cuenta.ObtenerSaldo()-saldoEsperado) < 0.001 {
		fmt.Printf("Then: Integridad validada. El saldo final es %.2f (las operaciones fallidas no alteraron el saldo).\n", cuenta.ObtenerSaldo())
	} else {
		fmt.Println("ERROR: Se rompió la integridad del saldo.")
	}

	jsonLog, err := cuenta.ExportarLogJSON("transacciones.json")
	if err != nil {
		fmt.Printf("ERROR al generar JSON: %v\n\n", err)
		return
	}

	fmt.Println("Resultado: ÉXITO. Log guardado en 'transacciones.json':")
	fmt.Println(jsonLog)
	fmt.Println()
}

// ==========================================
// =========== FIN ACÁPITE 3 ================
// ==========================================
