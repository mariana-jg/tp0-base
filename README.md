# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación. Respetar las entradas de log planteadas en los ejercicios, pues son las que se chequean en cada uno de los tests.

La corrección personal tendrá en cuenta la calidad del código entregado y casos de error posibles, se manifiesten o no durante la ejecución del trabajo práctico. Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).

# Solución - Estudiante: Mariana Juarez Goldemberg - 108441

### Ejercicio N°1:
Para el cumplimiento de este ejercicio se crea el script de bash pedido `generar-compose.sh` el cual recibe los parámetros que se indicaron en la consigna, utilizando un script de Python (lo consideré una mejor alternativa que hacer todo en bash) `mi-generador.py` que se ejecuta utilizando esos parámetros.

Dentro del script de Python se encuentran 3 constantes definidas para la generación del docker-compose:
* HEADER: con el nombre del compose y y el header de la declaración de los servicios a levantar.
* SERVER_BLOCK: con las definiciones necesarias para crear un container de un servidor.
* NETWORK_BLOCK: con las definiciones para crear las redes utilizadas dentro del proyecto.

Luego, para poder crear el container de la cantidad de clientes especificada por parámetro, se creó una función `client_block(n)` que permite definir un cliente según su número identificador, devolviendo el string formado para poder escribirlo en el archivo.

Para cumplir con el objetivo del ejercicio, se abre el archivo (con `with open` para asegurar su cierre), se escriben las diferentes constantes con los bloques para formar el docker-compose, repitiendo la escritura del cliente con un ID desde 1 hasta la cantidad que llegue por parámetro + 1.

#### Ejecución

`./generar-compose.sh <ARCHIVO_SALIDA> <CANTIDAD_DE_CLIENTES>`

### Ejercicio N°2:
Para el cumplimiento de este ejercicio, dentro del script de Python `mi-generador.py`, se agregó como propiedad al servicio de cliente y servidor un `volumes`. Se configuró para que cada uno de los servicios tome su archivo de configuración (`config.yaml` para el cliente y `config.ini` para el servidor) como volume, para que la información de esos archivos se persista fuera del container. Además, se eliminaron las variables de entorno que tenían precedencia sobre los archivos de configuración. 

#### Ejecución

`./generar-compose.sh <ARCHIVO_SALIDA> <CANTIDAD_DE_CLIENTES>`

### Ejercicio N°3:
Para el cumplimiento de este ejercicio, se crea el script de bash `validar-echo-server.sh` pedido, dentro del mismo se define un mensaje para enviar y se obtiene el puerto y la IP del servidor utilizando grep y extrayendo los valores. 

Luego se lanza un contenedor efímero de busybox (imagen de linux minimalista) que ya trae incorporado netcat. Se conecta a la red interna `tp0_testing_net` y dentro del contenedor se ejecuta `"echo '$message' | nc $SERVER_IP $SERVER_PORT"`, con este comando se genera el mensaje y se manda a netcat, abriendo una conexión con el servidor y enviándole ese mensaje. `nc` devuelve lo que el servidor responda, guardándolo en la variable answer. Por último, se realiza la verificación de que el server devolvió el mismo mensaje.

### Ejercicio N°4:
Para el cumplimiento de este ejercicio, se modificaron los sistemas de cliente y servidor para que se logre un graceful shutdown al recibir `SIGTERM`. En la consigna, se recomienda investigar sobre el flag `-t` del comando `docker compose down`. Ese flag determina la espera en segundos que se da para que el proceso termine por su cuenta, si después de ese tiempo sigue vivo, Docker envía `SIGKILL`, matándolo inmediatamente.

#### Servidor
Se añade un nuevo flag `_running`, reemplazando el `while true` que controlaba el bucle principal. Ante la notificación de una señal SIGTERM (por ejemplo con `docker compose down -t <N>`), se llama a la función `shutdown(self, signum, frame)`, esta función pone en false el `running` y cierra el socket de escucha del server. Además, cierra todas las conexiones de clientes que sigan abiertas.

Además, se configura `accept()` con `settimeout(1)` para que el servidor despierte periódicamente, detecte el shutdown y termine dentro del *tiempo de gracia* -t antes de que Docker envíe SIGKILL.

#### Cliente
Dentro del struct del cliente se define un canal de señales interno `done` que se utiliza para notificar el apagado del cliente. En la inicialización, se define un canal de señales del SO `signalChan`, que se le notificará cuando se envíe un SIGTERM, allí se define una go routine que corre en paralelo al bucle principal y espera un `SIGTERM` en `signalChan`. Cuando la señal llega, cierra el canal `done`, despertando al loop para que termine de forma ordenada.

En la función principal del cliente, tenemos dos funciones que chequean el shutdown. La primera, `mustStop`, chequea antes de crear el socket del cliente y la segunda, `awaitShutdown`, asegura que el cliente no quede dormido si llega la señal (el select elige lo que pase primero: `done` se cierra o pasó el tiempo).

### Ejercicio N°5:

Para el cumplimiento de este ejercicio, comencé con la definición de un protocolo de comunicación entre el cliente y el servidor en el directorio `/protocol` como `bet.go` . A continuación, describo la estructura de los datos que envía el cliente (serializados en big-endian) al servidor. Definí que el tamaño del paquete sea dinámico dada la naturaleza de la información.

* Length del payload (2B): framing - tamaño total del payload (menos estos 2 bytes).
* Identificador de agencia (1B): ID de la agencia (del cliente) que apuesta.
* Length del nombre (2B): largo del nombre de la persona que apuesta. [Máx. 65535 bytes]
* Nombre (Tamaño variable): nombre de la persona que apuesta.
* Length del apellido (2B): largo del apellido de la persona que apuesta. [Máx. 65535 bytes]
* Apellido (Tamaño variable): apellido de la persona que apuesta.
* Documento (8B): DNI de la persona que realiza la apuesta.
* Fecha de nacimiento (10B): fecha de nacimiento de la persona que apuesta. Con formato "YYYY-MM-DD".
* Número (2B): número apostado. 

#### Servidor

Se implementó el decodificador del lado del servidor una vez recibido el paquete dentro del archivo `protocol_codec.py`. Se leen los 2 bytes que identifican el largo del payload, luego exactamente se lee ese tamaño y se deserializa.
Además, como pedía la consigna, se implementaron las funciones para evitar los short writes (para garantizar que todos los bytes se envíen) y short reads (para garantizar que se lea por completo el paquete). 

##### Flujo del servidor
* Registra el socket entrante.
* Intenta leer y deserializar una apuesta del cliente.
* Almacena la apuesta.
* Envía un ACK (1B) al cliente (un "1").
* Cierra la conexión.

#### Cliente

Se implementó la lógica de serialización y envío de la apuesta en el cliente. Cada cliente construye una estructura `Bet` con los datos provenientes de las variables de entorno proporcionadas, y luego utiliza el método `ToBytes()` definido en el módulo `protocol` para serializarlo.
También se encuentran definidas las funciones para evitar los short writes y short reads.

##### Flujo del cliente
* Construye la apuesta a partir de las variables de entorno.
* Serializa la apuesta, devuelve el paquete en formato binario `[Len(2B) + payload]`.
* Abre un socket hacia el servidor.
* Envía la apuesta al servidor.
* Espera la confirmación del ACK.
* Loguea el resultado de la operación y cierra la conexión con el servidor.

#### Intento de reconexión en el cliente
Se agregó dentro del cliente luego de consultar en clase el reintento de conexión. Intenta conectarse hasta una cantidad `MAX_TRIES` definida como una constante. Decidí hacerlo de esta manera y no setear el reintento desde el Docker Compose para que el cliente pueda reconectarse ante cortes transitorios sin matar el proceso.

### Ejercicio N°6:

Para el cumplimiento de este ejercicio, realicé los siguientes agregados a la implementación:
* Lectura de las apuestas en archivo .CSV.
* Armado de los batches y serialización para que envíe el cliente.
* Deserialización de los batches en el servidor.

Para estos últimos dos puntos, se reutilizaron las funciones anteriores para las apuestas individuales.

Además, un detalle a destacar, es el cambio en la implementación del protocolo ya que ahora lo que se envía no es solamente una apuesta individual, sino el batch como tira de bytes.

Se abre un único socket hacia el servidor, se envían todas las apuestas y cuando se termina, se envía un mensaje de fin para cerrar la conexión.

### Ejercicio N°7:

Para el cumplimiento de este ejercicio, el cambio principal fue dentro del protocolo de comunicación, ya que se agrego como primer componente de los paquetes un byte que identificaba el tipo:
* TYPE_BET: para identificar a una apuesta que envía un cliente.
* TYPE_DONE: para identificar que el cliente terminó de realizar todas las apuestas.

Quedando el paquete enviado con esta estructura `[Type(1B) + Len(2B) + payload]`.

El servidor cuando le llega un paquete, lee el tipo del mismo y a partir de allí:
* Si es una apuesta, la procesa y envía un paquete de confirmación.
* Si es un mensaje del cliente avisando que terminó, procede a cortar el bucle de lectura y verifica si la cantidad de clientes que terminaron era la esperada que se conectaran (se añadió que el server conozca la cantidad de clientes que se conectarán, cambio realizado en el `main.py` y en `mi-generador.py`). Si todos los clientes ya avisaron que terminaron y se encuentran esperando por el resultado del sorteo, el sorteo se realiza y se envían los resultados de los ganadores a los clientes, cerrando luego sus respectivas conexiones.

### Ejercicio N°8:

Para el cumplimiento de este ejercicio, se utilizó la librería `multiprocessing` de Python utilizando `Process`, `Manager` y `Barrier` para evitar las limitaciones del GIL del lenguaje, trabajando con procesos independientes para manejar cada cliente.

Se utilizan barreras (en específico 2) que bloquean los procesos hasta que hayan llegado exactamente los clientes esperados a ese mismo punto.

##### Flujo del servidor concurrente

* Para cada cliente se crea un proceso que corre en segundo plano.
* Cada proceso se encarga de recibir y procesar todas las apuestas de su cliente. Cuando se recibe un paquete del tipo TYPE_DONE, se corta el bucle y pasa a sincronizarse.
* Tenemos una primer barrera para asegurar que todos hayan terminado de enviar las apuestas, algún proceso que llega (es random la elección) es el encargado de realizar el sorteo.
* Se utiliza una segunda barrera para que todos esperen a que el sorteo se realice (se cargan los ganadores en la estructura `Manager().dict()`), cada hijo responde a su cliente con su lista de ganadores y cierra. De esta manera, cuando salen de la barrera todos tienen el resultado listo.

Además, se utiliza un lock sobre `store_bets` para el manejo concurrente del recurso compartido, asegurando la exclusión mutua.

Se agregaron las modificaciones pertinentes para que el cliente al recibir el byte 255, lo tome como una caída del servidor y cierre su ejecución de forma segura.

