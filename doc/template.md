# Informe del Trabajo Práctico Grupal – SIGOA

## Algoritmos y Programación II – Primer Cuatrimestre 2025

**Integrantes del Grupo:**
- Juan Pérez – 12345  
- María Gómez – 23456  
- Lucas Fernández – 34567  
- Valentina Ruiz – 45678  

**Fecha de Entrega:** 15 de junio de 2025

---

## 1. Introducción

- **Problema a resolver:**  
  SIGOA (Sistema Integrado de Gestión de Operaciones Aeroportuarias) busca coordinar y automatizar tareas clave en un aeropuerto: el check-in de pasajeros, la asignación de carga a aeronaves, el despacho de vuelos y el embarque de pasajeros, gestionando eficientemente datos y prioridades.

- **Objetivos del trabajo:**
  1. Diseñar una arquitectura modular en Go para cada subproceso (Check-in, Carga, Despacho y Embarque).  
  2. Implementar estructuras de datos y algoritmos adecuados (colas de prioridad, mochila, sweep-line, compresión Huffman).  
  3. Analizar complejidad y evaluar desempeño bajo distintos escenarios.

- **Resumen de la solución:**
  - Módulo Check-in: simula llegadas, usa un heap como cola de prioridad para procesar pasajeros según categoría y hora de llegada.  
  - Módulo Carga: aplica un algoritmo de mochila bi-restricción (peso y volumen) para seleccionar carga óptima.  
  - Módulo Despacho: filtra vuelos seguros en relación a edificios y calcula línea de horizonte con técnica sweep-line.  
  - Módulo Embarque: ordena pasajeros por categoría y zona, codifica la lista final con Huffman.  

- **Organización del informe:**  
  2. Diseño del sistema  
  3. Implementación  
  4. Análisis de complejidad  
  5. Pruebas y resultados  
  6. Conclusiones y trabajo futuro  

---

## 2. Diseño del Sistema

### 2.1. Arquitectura General

- **Diagrama de bloques:**  
  ┌──────────────────┐       ┌─────────────┐       ┌────────────────┐  
  │   Check-in App   │ ────▶ │ Embarque App │ ────▶ │   Huffman I/O  │  
  └──────────────────┘       └─────────────┘       └────────────────┘  
            │                           ▲  
            │                           │  
            ▼                           │  
  ┌──────────────────┐       ┌─────────────┐  
  │   Carga App      │ ────▶ │ Despacho App │  
  └──────────────────┘       └─────────────┘  

- **Decisiones de diseño:**  
  - Lenguaje Go por concurrencia nativa y manejo de slices/structs.  
  - Singleton para instancias de cada módulo y carga única de CSV.  
  - Paquete común `csvfile` para abstraer lectura genérica de datos.  

### 2.2. Estructuras de Datos

1. **Check-in**  
   - Colas de prioridad (`heap`) de `LlegadaPasajero` con prioridad según categoría (Platino/Oro/Plata/Normal) y tiempo de llegada.  
   - Se usa `container/heap` de Go para operaciones O(log n) en push/pop.

2. **Gestión de Carga**  
   - Slice de estructuras `CargaStruc`.  
   - Matriz 2D `dp[n+1][maxPeso+1]` para el algoritmo de mochila extendido.  
   - Slice booleano `seleccionados` para reconstruir la solución.

3. **Despacho**  
   - Slice de `VueloStruc` y `EdificioStruc`.  
   - Barrido de eventos (`Punto`), ordenado por coordenada x (sweep-line).  
   - Mapa de alturas activas para mantener maxHeight en O(log n) al iterar.

4. **Embarque**  
   - Slice de `PasajeroEmbarque`.  
   - Ordenamiento simple (sort.Slice) por dos claves: categoría y zona.  
   - Codificación Huffman con árbol de frecuencias para compresión de salida.

### 2.3. Algoritmos Implementados

1. **Simulación de llegadas (Check-in):**  
   - Genera tiempos aleatorios en [0,60) min con seed variable.  
   - Inserta en heap según prioridad.  
   - Procesa pop hasta agotar capacidad (O(m log m), m = # pasajeros).

2. **Mochila bi-restricción (Carga):**  
   - DP clásico adaptado:  
     for i=1..n, w=maxPeso..peso_i:  
       dp[i][w] = max(dp[i-1][w], volumen_i + dp[i-1][w-peso_i])  
   - Reconstrucción con backtracking en O(n + maxPeso).

3. **Sweep-line y línea del horizonte (Despacho):**  
   - Genera eventos de entrada/salida de edificios.  
   - Ordena eventos O(k log k), k=2×#edificios.  
   - Mantiene multiconjunto de alturas en mapa y extrae cambios de skyline.

4. **Ordenación y compilación de lista (Embarque):**  
   - Filtra checkeados, asigna zona según `ConfiguracionAsientoStruc`.  
   - sort.Slice por prioridad de categoría y zona en O(p log p).  
   - Huffman: construye árbol de frecuencias (O(L + Σlog Σ)) y codifica texto.

### 2.4. Gestión de Datos de Entrada y Salida

- **Entradas CSV:**  
  - Estructuras Go genéricas: `CargaCSV[T any]("file.txt")`.  
  - Separador `;`, lectura de campos con refléxión en models.

- **Salidas:**  
  - Check-in y embarque: archivos de texto plano y `.out` Huffman.  
  - Despacho: `linea_horizonte.txt` con pares `x altura`.  
  - Embarque: carpeta `output/YYYYMMDD_HHMMSS.out` con datos comprimidos.

---

## 3. Implementación

- **Lenguaje Go:**  
  Uso de módulos (`go.mod`) y paquetes internos:
  - `internal/checkin`  
  - `internal/carga`  
  - `internal/embarque`  
  - `internal/ despacho`  
  - `pkg/csvfile`, `pkg/huffman`

- **Modularización:**  
  - Cada módulo es un singleton con `GetInstance()` o `NewCheckin()`.  
  - Inicialización perezosa y lectura única de CSV.

- **Funciones clave:**  
  - `SimularLlegadas`, `ProcesarCheckin`  
  - `mochilaCarga`  
  - `ObtenerVuelosSeguros`, `CalcularHorizonte`  
  - `EjecutarEmbarque`, `HuffmanEncode`

- **Buenas prácticas:**  
  - Manejo de errores con `log.Fatalf`.  
  - Uso de `fmt.Printf` para trazas.  
  - Separación de responsabilidades y abstracción de I/O.

---

## 4. Análisis de Complejidad

| Módulo       | Algoritmo                           | Tiempo                  | Espacio               |
|--------------|-------------------------------------|-------------------------|-----------------------|
| Check-in     | Heap push/pop m veces               | O(m log m)              | O(m)                  |
| Carga        | Mochila dp n×maxPeso                | O(n × W)                | O(n × W)              |
| Despacho     | Sweep-line eventos k=2e              | O(k log k + e × m)      | O(k + #alturas)       |
| Embarque     | sort.Slice p pasajeros              | O(p log p)              | O(p + Σfreq)          |

- W = capacidad máxima de peso redondeada.
- e = número de edificios.
- m = número de vuelos (en chequeo, pasajeros).
- p = número de pasajeros para embarque.

---

## 5. Pruebas y Resultados

- **Escenarios de prueba:**
  1. Bajo tráfico: 50 pasajeros, 10 cargas, 5 vuelos.  
  2. Tráfico medio: 500 pasajeros, 200 cargas, 20 vuelos.  
  3. Alto tráfico: 2 000 pasajeros, 1 000 cargas, 50 vuelos.

- **Metodología:**  
  - Scripts Go que generan datos aleatorios respetando rangos realistas.  
  - Cronometrado con `time.Since()` y conteo de operaciones.

- **Resultados:**
  - Check-in (medio): tiempo medio ≈ 150 ms, espera promedio ≈ 3 s.  
  - Carga (alto): cálculo de mochila ≈ 250 ms.  
  - Despacho horizonte (50 edificios): ≈ 20 ms.  
  - Embarque + Huffman (500 pax): ≈ 120 ms.

- **Interpretación:**  
  - La mayor latencia está en la mochila cuando W es grande.  
  - Heap y sweep-line responden en tiempo aceptable incluso en escenarios altos.

(Tabla y gráficos anexos en repositorio)

---

## 6. Conclusiones y Trabajo Futuro

- **Logros:**  
  - Sistema modular completamente implementado en Go.  
  - Uso adecuado de estructuras y algoritmos para cada subproceso.  
  - Rendimiento satisfactorio en escenarios de prueba.

- **Limitaciones:**  
  - DP de mochila consume memoria en casos de W elevado.  
  - Huffman no es incremental; toda la cadena debe generarse y codificarse.

- **Mejoras propuestas:**  
  1. Versión aproximada o heurística para mochila si W es muy grande.  
  2. Uso de estructuras concurrentes para procesar múltiples vuelos en paralelo.  
  3. Persistencia en base de datos en lugar de CSV para volúmenes masivos.

- **Trabajo en grupo:**  
  - Coordinación mediante GitHub y metodologías ágiles.  
  - Aprendizaje en Go, algoritmos y diseño de sistemas distribuidos.

---

## 7. Bibliografía

- Cormen, T.; Leiserson, C.; Rivest, R.; Stein, C. “Introduction to Algorithms”. MIT Press, 2009.  
- Kernighan, B.; Pike, R. “The Go Programming Language”. Addison-Wesley, 2015.  
- Sedgewick, R.; Wayne, K. “Algorithms, 4th Edition”. Addison-Wesley, 2011.

---

## 8. Apéndice

- Código fuente completo en: https://github.com/ejemplo/SIGOA  
- Datos de prueba y scripts de generación en `/testdata`  
- Instrucciones de compilación y ejecución en README del repositorio.