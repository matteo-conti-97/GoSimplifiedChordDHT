É stata implementata una versione semplificata di chord nella quale ogni nodo conosce solo i suoi diretti successori e predecessori,
il routing per individuare il peer responsabile di una risorsa avviene circolarmente, in particolare, se un peer non é responsabile per una risorsa inoltra la richiesta al suo successore
se é responsabile risponde al diretto richiedente (cioé viene percorsa a ritroso la catena di richieste). Il routing per cercare il peer responsabile una risorsa
fa ottenere al richiedente le coordinate di rete di tale peer, il richiedente poi esegue la PUT o la GET della risorsa contattando direttamente il peer responsabile.

Le risorse nella DHT sono state modellade come stringhe il cui Id é l'hash della stringa stessa (md5), gli Id dei nodi sono invece l'hash della stringa IP+porta del nodo
