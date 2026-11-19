# prj_withdrawal — Servicio ATM Web (UNI - FIIS)

Aplicación web en arquitectura de 3 capas (Go + HTML/CSS/JS + JSON) correspondiente al Acápite 4 del curso **SW608 - Lenguaje de Dominio Específico**.

## Funcionalidades
- **US-01a / RN-01, RN-02:** Retiro de efectivo con validación de saldo y montos negativos.
- **US-01b:** Abono en efectivo.
- **US-04 / RN-03:** Integridad atómica de transacciones y registro de auditoría en `transacciones.json` (`EXITO`, `FALLO`, `ERROR`).

## Ejecución
```bash
# Ejecutar pruebas BDD en consola
go run main.go -bdd

# Iniciar la App Web (http://localhost:8080)
go run main.go
