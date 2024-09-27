# counterMag - A distributed, in-memory database specialized in occurence counting

## Estimari

- Cereri analiza text: 1 000 req / s
- Cereri citire text:  10 000 req / s
- Dimensiune medie text: 100 cuvinte

## Cerinte functionale

- Analizeaza un text si stocheaza numarul de aparitii pentru fiecare cuvant
- Obtine numarul de aparitii pentru un set de cuvinte date
- Datele sa persiste in urma repornirii aplicatiei

## Cerinte non-functionale

- Sa poata scala prin cresterea numarului de instante pentru a face fata workload-ului
- Consistenta eventuala
- Sistemul sa ramana functional daca una din instante devine indisponibila


## Persistenta

Persitenta bazei de date este impementata prin salvarea periodica a continutului acesteia intr-un fisier `counter-{instanceId}.log`, unde `instanceId` este ID-ul unic al instantei de baza de date.

Modul prin care se realizeaza acest lucru este la nivelul clasei [`Database`](/internal/database/database.go#L77). Pentru a persista continutul, clasa `Database` foloseste un struct care implementeaza interfata [`SnapshotPersister`](/internal/database/snapshot.go#L11). Aceasta interfata expune metodele `SaveSnapshot` si `LoadSnapshot` care sunt folosite pentru scrierea, respectiv citirea unui snapshot intr-un mod definit de orice `struct` care implementeaza interfata respectiva.

In practica, struct-ul concret folosit este [`FileSnapshotPersister`](/internal/database/snapshot.go#L16), care salveaza cheile din baza de date in formatul `<cheie> <valoare>`, separate de caracterul `\n`

Aceasta interfata a aparut din nevoia de testare a clasei `Database` fara a interactiona in mod direct cu file system-ul, in teste fiind folosit [`MockSnapshotPersister`](/internal/database/snapshot.go#L71) care doar tine datele nemodificate intr-un camp.

Datele sunt salvate periodic folosind un goroutine separat pornit prin apelul functiei [`StartPeriodicFlush`](/internal/database/database.go#L90). Aceasta este apelata implicit in functia constructor [`NewDatabase`](/internal/database/database.go#L18).

Acest goroutine este oprit in momentul in care contextul folosit de `Database` este 'oprit'.

`Database` exporta si functia [`Close`](/internal/database/database.go#L106) care este apelata la momentul opririi aplicatiei cu un [timeout](/cmd/countermag/main.go#L78) ca parte din procesul de graceful shutdown pentru a evita pierderea datelor.


## Replicare

Aplicatia expune doua servere HTTP:
- cel de aplicatie, responsabil de analiza textelor si interogarile pentru aparitiile de cuvinte
- cel de cluster, responsabil de facilitatea comunicarii intre instantele de baza de date

Modelul de replicare ales este cel de `master-slave`, dupa cum urmeaza:
- fiecare instanta este pornita cu flag-urile `cluster` si `port`, unde `cluster` reprezinta adresa nodului master in format `host:port` si `port` reprezinta portul pe care va rula nodul curent
- daca portul nodului curent coincide cu portul din adresa nodului master, nodul curent va porni ca master, altfel se va conecta la master si acesta il va inregistra ca slave
- scrierile vor merge catre master si vor fi replicate periodic catre slaves
- citirile vor merge catre slaves

Este de mentionat ca rutarea scrierilor si citirilor este revizitata la sectiunea de imbunatatiri

Rutele expuse de serverul de cluster sunt urmatoarele:
- `GET    /cluster` - apelat pe nodul master va intoarce toate nodurile conectate la cluster alaturi de portul si starea lor
- `PUT  /store` - apelat periodic de catre master pentru fiecare dintre nodurile slave pentru a sincroniza baza de date
- `POST /connect` - apelat de catre slaves catre master pentru a se alatura clusterului
- `GET  /ping` - folosit de master pentru a verifica periodic starea nodurilor slave

Logica pentru apelarea acestor endpoint-uri a fost incapsulata in [`ClusterClient`](/internal/cluster/client.go#L11)

