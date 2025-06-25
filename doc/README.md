# Informe del Trabajo Práctico Grupal - SIGOA

## Algoritmos y Programación II - Primer Cuatrimestre 2025

**Integrantes del Grupo:**

  - Andres Lacquaniti

**Fecha de Entrega:** 24 de junio de 2025

## 1\. Introducción

El presente informe detalla el desarrollo de un **Sistema Integrado de Gestión de Operaciones Aeroportuarias (SIGOA)**, cuyo objetivo principal es simular y optimizar los procesos de un aeropuerto, incluyendo el check-in de pasajeros, la gestión de carga y el despacho de vuelos. Este proyecto busca aplicar los conceptos de estructuras de datos y algoritmos avanzados para diseñar e implementar soluciones eficientes.

La solución implementada se enfoca en gestionar las operaciones aeroportuarias utilizando diversas estructuras de datos como listas enlazadas y colas de prioridad para el manejo de pasajeros, y algoritmos de optimización para la carga y verificación de seguridad de vuelos.

El informe se organiza en las siguientes secciones: diseño del sistema, estructuras de datos y algoritmos implementados, complejidad, pruebas y resultados, y conclusiones y trabajo futuro.

## 2. Diseño del Sistema

### 2.1. Arquitectura General

El sistema SIGOA está modularizado en varios paquetes Go, cada uno encargado de una funcionalidad específica, lo que facilita la mantenibilidad y escalabilidad. La arquitectura se puede describir mediante los siguientes módulos principales y su interacción:

  * **`app`**: Maneja la lógica central de la simulación del tiempo del sistema (`HoraSistema`) y la coordinación general.
  * **`models`**: Define las estructuras de datos (`structs`) que representan las entidades del dominio, como `AeronaveStruc`, `AeropuertoStruc`, `CargaStruc`, `ClienteStruc`, `ConfiguracionAsientoStruc`, `EdificioStruc`, `ReservaStruc`, `VueloStruc`, y `LlegadaPasajero`.
  * **`utils`**: Provee funciones de utilidad transversal, como el registro (`PrintLog`) y la generación de números de ticket (`GeneraNroTicket`).
  * **`csvfile` (pkg)**: Encargado de la carga de datos desde archivos CSV a las estructuras Go correspondientes.
  * **`checkin`**: Gestiona todo el proceso de check-in de pasajeros, incluyendo la simulación de llegadas, la atención en mostradores y la gestión de colas de prioridad.
  * **`vuelo`**: Administra la información y el estado de los vuelos, incluyendo la carga de datos de vuelos y edificios, y la verificación de seguridad de los mismos.
  * **`carga`**: Se ocupa de la optimización y procesamiento de la carga para cada vuelo.
  * **`embarque`**: Coordina el proceso de embarque de pasajeros y la carga de los vuelos.
  * **`guardar`**: Responsable de persistir los registros finales de los vuelos en formato JSON y comprimirlos en huffman.
  * **`huffman` (pkg)**: Utilizado para la compresión de datos de los registros finales.


### 2.1.1 Diagrama de Bloques del Sistema
Este diagrama ilustra la interacción entre los principales módulos del sistema:
```
+------------------+      +---------------------+
|      app         |<-----|       utils         |
| (Tiempo Simulado)|      |  (Logs, Tickets)    |
+------------------+      +---------------------+
        |                          ^
        |                          |
        v                          |
+-----------------+      +---------------------+
|     csvfile     |----->|       models        |
| (Carga de Datos)|      |   (Estructuras)     |
+-----------------+      +---------------------+
        |
        v
+---------------------+    +---------------------+     +------------------+
|     checkin         |--->|        vuelo        |<----|    horizonte     |
| (Gestión Pasajeros) |    |   (Estado Vuelos,   |     |  (Cálculo Altura)|
| (Colas de Prioridad)|    |     Seguridad)      |     +------------------+
+---------------------+    +---------------------+
        |                          |
        v                          v
+--------------------+      +---------------------+
|    Mostrador       |----->|        carga        |
| (Atención Check-in)|      |   (Optimización     |
+-----------------+         |      de Carga)      |
        |                   +---------------------+
        v                            ^
+-----------------+                  |
|    embarque     |------------------+
| (Coordinación   |
|   Embarque)     |
+-----------------+
        |
        v
+-----------------+
|     guardar     |
| (Registro Final |
|    y JSON)      |
+-----------------+
```

### 2.1.2 Diagramas de Apoyo

**Modelo Entidad-Relación (DER):**


![der](./der.drawio.svg)

**Flujo General del Proyecto:**


![der](./workflow.drawio.svg)



La arquitectura adoptada es modular y orientada a la separación de responsabilidades, lo que permite que cada componente se especialice en una tarea y se comunique con los demás a través de interfaces bien definidas. La justificación de esta decisión radica en la mejora de la cohesión del código y la reducción del acoplamiento entre los módulos, facilitando así el desarrollo, las pruebas y el mantenimiento del sistema.

### 2.2. Estructuras de Datos

El sistema hace un uso extensivo de diversas estructuras de datos para modelar las entidades y gestionar los flujos de trabajo.

#### Módulo `models` (General)

Este paquete define todas las estructuras de datos base utilizadas en el sistema para mapear los datos de entrada CSV y la información interna.

  * **`AeronaveStruc`**: Representa una aeronave con su matrícula, número de asientos, capacidad de peso y volumen de carga.
  * **`AeropuertoStruc`**: Detalla un aeropuerto con su provincia, ciudad, nombre y código IATA.
  * **`CargaStruc`**: Define un ítem de carga con su destino, peso y volumen.
  * **`ClienteStruc`**: Almacena la información de un cliente (pasajero) incluyendo nombre, apellido, DNI y categoría. Incluye etiquetas CSV y JSON.
  * **`ConfiguracionAsientoStruc`**: Guarda la configuración de asientos por aeronave y zona.
  * **`EdificioStruc`**: Representa un edificio con sus coordenadas `xi`, `altura` y `xf`, utilizado para el cálculo del horizonte.
  * **`ReservaStruc`**: Contiene los datos de una reserva, como el código de reserva, DNI del pasajero, número de vuelo, fecha y estado.
  * **`VueloStruc`**: Representa un vuelo con su número, fecha y hora, destino y matrícula de la aeronave.
  * **`LlegadaPasajero`**: Estructura clave para el check-in, que encapsula el ticket, DNI, prioridad de embarque, hora de llegada, y zonas asignadas. También incluye un `Index` para la cola de prioridad.

#### Módulo `checkin`

  * **`QueueLlegada` (tipo `*list.List`)**: Utilizada como una cola FIFO (First-In, First-Out) para simular la llegada de pasajeros. Los pasajeros que llegan son encolados aquí antes de ser procesados por los mostradores.
      * **Justificación**: Una lista enlazada es adecuada para la simulación de llegadas, ya que permite inserciones rápidas al final y eliminaciones rápidas al principio, emulando el flujo de una cola real.
  * **`PrioridadQueue` (tipo `[]*models.LlegadaPasajero` que implementa `heap.Interface`)**: Es una cola de prioridad (min-heap) que organiza a los pasajeros para el embarque. La prioridad se determina por la `Prioridad` del pasajero y, en caso de empate, por la hora de llegada (`Llegada`).
      * **Justificación**: Una cola de prioridad es fundamental para asegurar que los pasajeros con mayor prioridad (por su categoría) sean atendidos primero, y entre aquellos con la misma prioridad, se respete el orden de llegada. La implementación de la interfaz `heap` de Go proporciona eficiencia para las operaciones de `Push`, `Pop` y `Actualizar`.
  * **`Pqueue` (tipo `map[string]*PrioridadQueue`)**: Un mapa que almacena una `PrioridadQueue` por cada número de vuelo. Esto permite gestionar colas de embarque separadas para cada vuelo.
      * **Justificación**: Un mapa es ideal para acceder rápidamente a la cola de prioridad de un vuelo específico utilizando su número como clave.
  * **`Reservas` (tipo `[]models.ReservaStruc`)**: Slice que almacena todas las reservas cargadas desde el archivo `reservas.txt`.
  * **`Clientes` (tipo `[]models.ClienteStruc`)**: Slice que contiene todos los clientes cargados desde `clientes.txt`.
  * **`categoriaPrioridad` (tipo `map[string]int`)**: Mapa estático que asigna un valor numérico de prioridad a cada categoría de cliente (Platino, Oro, Plata, Normal).
  * **`ZonaSalida` (tipo `map[string]int`)**: Mapa estático que asigna una zona de salida a cada categoría de cliente.

#### Módulo `vuelo`

  * **`VueloApp`**: Estructura que representa la instancia de un vuelo con su información (`models.VueloStruc`) y su `Estado` actual (Cerrado, CheckIn, Embarque, PreDespegue, Despegue).
  * **`singleInstance` (tipo `map[string]*VueloApp`)**: Un mapa que implementa el patrón Singleton para las instancias de `VueloApp`, asegurando que solo exista una instancia de `VueloApp` para cada número de vuelo.
      * **Justificación**: Permite gestionar de forma centralizada y eficiente el estado y las operaciones de cada vuelo sin duplicar recursos.
  * **`vuelos` (tipo `[]models.VueloStruc`)**: Cache de todos los vuelos cargados, ordenados por fecha y hora de salida.
  * **`Edificios` (tipo `[]models.EdificioStruc`)**: Slice que contiene la información de todos los edificios cargados, utilizada para la verificación de seguridad de vuelos.

#### Módulo `carga`

  * **`CargaApp`**: Estructura principal que gestiona las aeronaves y cargas.
  * **`instancia` (tipo `*CargaApp`)**: Puntero para implementar el patrón Singleton para la aplicación de carga.
      * **Justificación**: Similar al módulo `vuelo`, asegura una única gestión de los recursos de carga.
  * **`Aeronaves` (tipo `[]models.AeronaveStruc`)**: Slice que almacena todas las aeronaves disponibles.
  * **`Cargas` (tipo `[]models.CargaStruc`)**: Slice que almacena todas las cargas disponibles.

#### Módulo `guardar`

  * **`RegistroFinalStruc` (tipo `map[string]*RegistroVueloStruc`)**: Un mapa que almacena una estructura `RegistroVueloStruc` para cada vuelo, registrando los pasajeros embarcados, no presentes y en lista de espera.
      * **Justificación**: Permite consolidar toda la información relevante de cada vuelo para su posterior guardado y análisis.
  * **`RegistroVueloStruc`**: Estructura que consolida los datos finales de un vuelo, incluyendo el número de vuelo, hora programada y real de salida, listas de pasajeros embarcados, no presentes y en espera.

### 2.3. Algoritmos Implementados

El sistema integra varios algoritmos clave para sus funcionalidades principales:

#### Módulo `app`

  * **Simulación de Tiempo (`Init` función en `app/main.go`)**: Un `goroutine` se encarga de incrementar una variable `HoraSistema` a intervalos regulares, simulando el paso del tiempo real en el sistema.
      * **Tarea clave**: Proporcionar una línea de tiempo simulada para la ejecución de los eventos del aeropuerto, desacoplando la simulación del tiempo real.

#### Módulo `checkin`

  * **`ObtenerPasajerosPorVuelo`**: Itera sobre las reservas y clientes para encontrar los pasajeros confirmados para un vuelo específico.
      * **Tarea clave**: Filtrar y agrupar pasajeros por vuelo.
  * **`SimularLlegadas`**: Introduce a los pasajeros en la `QueueLlegada` con una pausa aleatoria para simular llegadas asincrónicas.
      * **Tarea clave**: Simular la afluencia de pasajeros al aeropuerto.
  * **`CalculaMostradores`**: Determina dinámicamente el número de mostradores a abrir en función de la longitud de la cola de llegadas.
      * **Tarea clave**: Asignar recursos (mostradores) según la demanda.
  * **`StartMostrador`**: Utiliza `goroutines` (uno por mostrador) para procesar pasajeros de la `QueueLlegada` de forma concurrente, asegurando que se esperen a que todos los mostradores asignados terminen de procesar un lote de pasajeros antes de proceder con el siguiente.
      * **Tarea clave**: Procesar pasajeros en mostradores de check-in de forma concurrente.
  * **`Mostrador`**: Simula la atención individual de un pasajero en un mostrador, calculando el tiempo de espera.
      * **Tarea clave**: Representar el proceso de atención en un mostrador.
  * **`ProcesarCheckin`**: Gestiona el proceso de check-in de un pasajero, verificando si llegó tarde para el vuelo y agregándolo a la `Pqueue` (cola de prioridad) del vuelo si corresponde, o a la lista de "no presentes".
      * **Tarea clave**: Finalizar el check-in, asignando prioridades y zonas, y gestionando pasajeros tardíos.
  * **Cola de Prioridad (`PrioridadQueue`)**: La implementación de la interfaz `heap.Interface` permite operaciones eficientes de inserción (`Push`), eliminación (`Pop`) y actualización (`Actualizar`) de pasajeros basándose en su prioridad y hora de llegada.
      * **Tarea clave**: Organizar y despachar pasajeros para el embarque según reglas de prioridad.

#### Módulo `vuelo`

  * **`ActualizarEstado`**: Determina el estado actual de un vuelo (Cerrado, CheckIn, Embarque, PreDespegue, Despegue) basándose en la `HoraSistema` simulada y la `FechaHora` de salida programada del vuelo.
      * **Tarea clave**: Gestionar el ciclo de vida del vuelo.
  * **`VuelosSeguro` (en `vuelo/horizonte.go`)**: Genera una altitud aleatoria y verifica si el vuelo potencial chocaría con algún edificio cargado.
      * **Tarea clave**: Asegurar la seguridad de los despegues evitando colisiones con edificios.
  * **`CalcularHorizonte` (en `vuelo/horizonte.go`)**: Implementa un algoritmo de *sweep line* para calcular la línea del horizonte a partir de un conjunto de edificios y guarda el resultado en un archivo.
      * **Tarea clave**: Determinar el perfil del horizonte para la seguridad de los vuelos.

#### Módulo `carga`

  * **`mochilaCarga`**: Implementa una variación del algoritmo de la mochila (Knapsack Problem) para seleccionar la combinación óptima de cargas que maximice el volumen utilizado sin exceder la capacidad de peso y volumen de la aeronave.
      * **Tarea clave**: Optimizar la carga de la aeronave.

#### Módulo `guardar`

  * **`GuardarRegistroVueloEnJson`**: Serializa el registro final de un vuelo a formato JSON y lo guarda en un archivo.
      * **Tarea clave**: Persistir los datos finales de la simulación de cada vuelo.
  * **`ComprimirArchivo` (usando `huffman`):** Comprime el archivo JSON generado utilizando el algoritmo de Huffman.
      * **Tarea clave**: Almacenamiento eficiente de los datos.

## 3\. Análisis de Complejidad

### 3.1. Complejidad Temporal y Espacial

#### **`checkin`**

  * **`ObtenerPasajerosPorVuelo`**:
      * **Temporal**: $O(R \\cdot C)$, donde $R$ es el número de reservas y $C$ es el número de clientes. En el peor caso, itera sobre todas las reservas y, para cada una, sobre todos los clientes.
      * **Espacial**: $O(P)$, donde $P$ es el número de pasajeros para el vuelo específico (el slice `pasajeros`).
  * **`SimularLlegadas`**:
      * **Temporal**: $O(P \\cdot T\_{\\text{sleep}})$, donde $P$ es el número de pasajeros y $T\_{\\text{sleep}}$ es el tiempo de pausa aleatoria entre llegadas. La función `rand.Intn` tiene una complejidad constante. La operación `PushBack` de `container/list` es $O(1)$.
      * **Espacial**: $O(P)$, para almacenar los objetos `LlegadaPasajero` en la cola.
  * **`StartMostrador`**:
      * **Temporal**: El bucle principal se ejecuta mientras la cola de llegadas no esté vacía. Dentro del bucle, se lanzan `nMostradores` goroutines. La función `Mostrador` tiene un `time.Sleep` variable. La complejidad total dependerá del número de pasajeros y del tiempo de procesamiento en cada mostrador. Sin embargo, la gestión de `sync.WaitGroup` es eficiente. En el peor de los casos, donde se procesa un pasajero por vez, sería $O(P \\cdot T\_{proc})$, donde $T\_{proc}$ es el tiempo de procesamiento de un pasajero.
      * **Espacial**: $O(1)$ adicional más allá de la cola de llegadas, ya que los datos se procesan y se eliminan de la cola.
  * **`PrioridadQueue` (Heap operations)**:
      * **`Len`, `Swap`**: $O(1)$.
      * **`Less`**: $O(1)$ (comparación de dos elementos).
      * **`Push` (en `heap`)**: $O(\\log N)$, donde $N$ es el número de elementos en la cola.
      * **`Pop` (en `heap`)**: $O(\\log N)$.
      * **`Actualizar` (en `heap.Fix`)**: $O(\\log N)$.
      * **Espacial (Heap)**: $O(N)$ para almacenar los elementos del heap.

#### **`vuelo`**

  * **`VuelosSeguro`**:
      * **Temporal**: $O(M)$, donde $M$ es la cantidad de edificios. Se itera sobre todos los edificios para verificar colisiones.
      * **Espacial**: $O(1)$ adicional, ya que solo se utilizan variables para la verificación.
  * **`CalcularHorizonte`**:
      * **Temporal**: $O(E \\log E)$, donde $E$ es el número de "puntos" (inicio y fin de edificios), que es $2N$ para $N$ edificios. El cuello de botella es el ordenamiento de los puntos. La iteración sobre los puntos y las operaciones en el mapa (`heights`) son en promedio $O(1)$ para inserciones y eliminaciones, pero en el peor caso (muchas colisiones de hash) podrían ser $O(H)$ donde $H$ es el número de alturas distintas activas. Sin embargo, dado el uso de un mapa para alturas, se mantiene eficiente.
      * **Espacial**: $O(E)$ para almacenar los puntos, y $O(H\_{max})$ para el mapa `heights` donde $H\_{max}$ es la altura máxima.

#### **`carga`**

  * **`mochilaCarga`**:
      * **Temporal**: $O(N \\cdot W\_{max})$, donde $N$ es el número de ítems de carga y $W\_{max}$ es la capacidad máxima de peso de la aeronave. Este es un algoritmo de programación dinámica.
      * **Espacial**: $O(N \\cdot W\_{max})$ para la matriz `dp`.

#### **`utils`**

  * **`PrintLog`**:
      * **Temporal**: Principalmente depende de las operaciones de E/S de archivo, que varían según el sistema operativo, pero conceptualmente es $O(L)$, donde $L$ es la longitud del mensaje. Las operaciones de creación/apertura de directorio y archivo son generalmente eficientes.
      * **Espacial**: $O(L)$ para el mensaje y la línea de log.
  * **`GeneraNroTicket`**:
      * **Temporal**: $O(K)$, donde $K$ es la longitud del hash (6 en este caso). La operación `rand.Int` es eficiente.
      * **Espacial**: $O(K)$ para el `strings.Builder`.

## 4\. Pruebas y Resultados

Para evaluar el rendimiento del sistema SIGOA, Se definieron 3 esenarios, se dejan en el directorio output

  * **Tráfico Bajo:** Pocos vuelos, pocos pasajeros por vuelo, pocas cargas.
  * **Tráfico Medio:** Número moderado de vuelos, pasajeros y cargas.
  * **Tráfico Alto:** Gran cantidad de vuelos, pasajeros y cargas, con posibles solapamientos de horarios.

### Metodología de las Pruebas

Las pruebas se realizarán ejecutando el programa principal con diferentes conjuntos de datos de entrada (`aeronaves.txt`, `aeropuertos.txt`, `cargas.txt`, `clientes.txt`, `configuracion_asientos.txt`, `edificios.txt`, `reservas.txt`, `vuelos.txt`) que simulan los diferentes escenarios de tráfico.

Se medirán las siguientes métricas:

  * **Tiempo promedio de espera en el check-in**: Calculado por `Mostrador` como la diferencia entre `HoraSistema` y `Llegada` del pasajero.
  * **Porcentaje de utilización de la capacidad de carga**: Se evaluará la eficiencia del algoritmo `mochilaCarga2` en el módulo `carga` comparando el volumen y peso cargado contra la capacidad total de la aeronave.
  * **Tiempo total de procesamiento de un vuelo**: Se registrará desde el inicio del check-in hasta el despegue simulado (cuando el estado del vuelo sea "Despegue").
  * **Número de pasajeros embarcados vs. no presentes**: Información registrada en `RegistroFinalStruc`.
  * **Número de conflictos de vuelo con edificios**: Contabilizados en el módulo `vuelo` por la función `VuelosSeguro`.
  * **Tamaño de los archivos de registro y JSON comprimidos**: Para evaluar la eficiencia de la compresión Huffman.

### Resultados Esperados (Ejemplos hipotéticos)

**Escenario de Tráfico Bajo:**

  * Tiempo promedio de espera en check-in: \~1-5 minutos.
  * Utilización de carga: \>90% (alta eficiencia debido a la menor cantidad de restricciones).
  * Tiempo total de procesamiento por vuelo: \~30-60 minutos (simulados).
  * Pasajeros no presentes: Bajo (pocos conflictos por llegada tardía).
  * Conflictos de vuelo: Muy bajo/nulo.
  * Tamaño de archivos: Pequeños, compresión efectiva.

**Escenario de Tráfico Medio:**

  * Tiempo promedio de espera en check-in: \~5-15 minutos.
  * Utilización de carga: \~80-90% (buena eficiencia).
  * Tiempo total de procesamiento por vuelo: \~60-120 minutos (simulados).
  * Pasajeros no presentes: Moderado.
  * Conflictos de vuelo: Ocasionales, resueltos por reintentos de altitud.
  * Tamaño de archivos: Medianos, buena compresión.

**Escenario de Tráfico Alto:**

  * Tiempo promedio de espera en check-in: \>15 minutos, con picos significativos.
  * Utilización de carga: \~70-85% (puede verse afectada por la diversidad de ítems).
  * Tiempo total de procesamiento por vuelo: \>120 minutos (simulados), con posibles demoras.
  * Pasajeros no presentes: Alto (mayor probabilidad de llegadas tardías).
  * Conflictos de vuelo: Más frecuentes, indicando la necesidad de múltiples reintentos para encontrar una altitud segura.
  * Tamaño de archivos: Grandes, la compresión será más relevante.

### Análisis de los Resultados

Pendiente de Realizacion.


## 5\. Conclusiones y Trabajo Futuro

### 5.1. Logros y Evaluación de Eficiencia

El desarrollo del Sistema Integrado de Gestión de Operaciones Aeroportuarias (SIGOA) ha permitido la aplicación práctica de diversos conceptos de estructuras de datos y algoritmos. Se ha logrado simular de manera efectiva procesos clave como el check-in de pasajeros, la gestión de carga y el despacho de vuelos.

  * **Modularidad**: La arquitectura basada en paquetes separados (`app`, `models`, `checkin`, `vuelo`, `carga`, `embarque`, `guardar`, `utils`) ha facilitado el desarrollo y la comprensión de las responsabilidades de cada componente.
  * **Gestión de Concurrencia**: El uso de `goroutines` y `sync.WaitGroup` en `StartMostrador` ha permitido simular la atención concurrente de pasajeros en mostradores, mostrando una mejora en el tiempo de procesamiento general.
  * **Optimización de Recursos**: La implementación de la cola de prioridad para el embarque garantiza que los pasajeros sean procesados según su categoría y hora de llegada, optimizando el flujo de embarque.
  * **Seguridad de Vuelos**: El cálculo del horizonte y la verificación de altitud para los vuelos (`VuelosSeguro` en `vuelo/horizonte.go`) es una característica crítica para la seguridad operativa.
  * **Eficiencia de Carga**: El algoritmo de la mochila para la carga de aeronaves (`mochilaCarga2`) aborda un problema de optimización conocido, buscando maximizar el uso de la capacidad disponible.
  * **Persistencia de Datos**: La capacidad de guardar registros detallados de los vuelos en formato JSON y comprimirlos es valiosa para análisis post-simulación.

### 5.2. Cuellos de Botella y Limitaciones

  * **Simulación de Tiempo Lineal**: El avance de `HoraSistema` es lineal y se basa en pausas fijas de tiempo real. Esto podría no escalar eficientemente para simulaciones muy largas o con un gran número de eventos que requieran saltos de tiempo significativos.
  * **I/O de Archivos**: Las operaciones de carga de CSV y escritura de logs y JSON pueden convertirse en cuellos de botella en escenarios de alto volumen de datos, especialmente si no se manejan de forma asincrónica o con buffering adecuado.
  * **Complejidad del Algoritmo de Mochila**: La complejidad $O(N \\cdot W\_{max})$ del algoritmo de la mochila puede ser un factor limitante si el número de ítems de carga o la capacidad de peso de las aeronaves son extremadamente grandes, aunque para casos realistas de carga de aeronaves, suele ser aceptable.
  * **Gestión de Concurrencia en Archivos de Log**: Aunque `PrintLog` es concurrencia-segura en cuanto a la apertura y escritura, el alto volumen de logs podría impactar el rendimiento.
  * **Aleatoriedad Fija**: El uso de `time.Now().UnixNano()` en `generarAltitudRandom` puede generar la misma secuencia de números aleatorios si se llama muy rápidamente. Aunque `rand.NewSource` con `UnixNano` se usa, para simulaciones muy precisas, podría ser preferible un generador de números pseudoaleatorios con una semilla de control más fina o un `crypto/rand` para mayor aleatoriedad.

### 5.3. Propuestas de Mejoras y Extensiones

  * **Simulación de Eventos Discretos**: Implementar un planificador de eventos para una simulación basada en eventos discretos en lugar de un tick de tiempo fijo. Esto permitiría saltar directamente entre eventos significativos, mejorando la eficiencia para simulaciones a gran escala.
  * **Base de Datos para Persistencia**: Reemplazar la lectura de CSV y escritura de JSON con una base de datos (por ejemplo, SQLite, PostgreSQL) para una gestión de datos más robusta, escalable y con capacidad de consulta.
  * **Interfaz de Usuario (Web/CLI)**: Desarrollar una interfaz de usuario básica (CLI avanzada o web) para visualizar el estado de la simulación en tiempo real, interactuar con el sistema (ej. agregar un vuelo, consultar estado de un pasajero) y presentar informes de forma más interactiva.
  * **Visualización del Horizonte**: Extender la función `CalcularHorizonte` para generar una representación gráfica (SVG, PNG) de la línea del horizonte.
  * **Algoritmos de Asignación de Recursos**: Implementar algoritmos más complejos para la asignación dinámica de asientos y puertas de embarque, posiblemente utilizando técnicas de *graph theory* o *constraint programming*.
  * **Optimización de Rutas de Vuelo**: Incluir un módulo para la planificación de rutas de vuelo óptimas, considerando factores como el consumo de combustible, condiciones climáticas y espacio aéreo.
  * **Monitoreo y Alertas**: Implementar un sistema de monitoreo que genere alertas ante situaciones críticas (ej. alta congestión en check-in, vuelos con demoras críticas).
