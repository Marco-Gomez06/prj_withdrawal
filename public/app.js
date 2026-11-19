const elTitular = document.getElementById("titular");
const elSaldo = document.getElementById("saldo");
const elMonto = document.getElementById("monto");
const elAlerta = document.getElementById("alerta");
const elTabla = document.getElementById("tabla-historial");
const elJson = document.getElementById("json-viewer");

function renderizarEstado(data, mostrarAlerta = false) {
  elTitular.textContent = data.titular;
  elSaldo.textContent = `S/ ${Number(data.saldo).toFixed(2)}`;

  if (mostrarAlerta && data.mensaje) {
    elAlerta.textContent = `[${data.estado}] ${data.mensaje}`;
    elAlerta.className = `alerta ${data.estado}`;
  }

  const historial = data.historial || [];
  if (historial.length === 0) {
    elTabla.innerHTML = `<tr><td colspan="6" class="vacio">Sin transacciones registradas aún.</td></tr>`;
  } else {
    elTabla.innerHTML = historial.map((tx) => `
      <tr>
        <td>${tx.fecha}</td>
        <td><strong>${tx.tipo}</strong></td>
        <td>S/ ${Number(tx.monto).toFixed(2)}</td>
        <td><span class="badge badge-${tx.estado}">${tx.estado}</span></td>
        <td>${tx.detalle}</td>
        <td>S/ ${Number(tx.saldo_final).toFixed(2)}</td>
      </tr>
    `).join("");
  }

  elJson.textContent = JSON.stringify(historial, null, 2);
}

async function cargarEstado() {
  const res = await fetch("/api/estado");
  const data = await res.json();
  renderizarEstado(data, false);
}

async function ejecutarOperacion(endpoint) {
  const monto = parseFloat(elMonto.value);
  if (isNaN(monto)) {
    elAlerta.textContent = "[ERROR] Ingrese un valor numérico válido.";
    elAlerta.className = "alerta ERROR";
    return;
  }

  const res = await fetch(endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ monto }),
  });
  const data = await res.json();
  renderizarEstado(data, true);
  elMonto.value = "";
}

document.getElementById("btn-retirar").addEventListener("click", () => {
  ejecutarOperacion("/api/retirar");
});

document.getElementById("btn-abonar").addEventListener("click", () => {
  ejecutarOperacion("/api/abonar");
});

document.getElementById("btn-reiniciar").addEventListener("click", async () => {
  const res = await fetch("/api/reiniciar", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ titular: "Marco Gomez", saldo_inicial: 200.00 }),
  });
  const data = await res.json();
  renderizarEstado(data, true);
});

cargarEstado();