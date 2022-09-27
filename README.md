# Software-Engineering-Chat0815-HWR
## branch for file transfer feature

### Open:
- IP-Zuweisung -> Laurin
- Anbindung an main -> Laurin
- GUI: Derzeit schließen die Fenster nicht automatisch und mann muss sie manuell schließen. Dies führ dazu dann man nicht mehrere dateien senden kann. Um mehrere
Dateien zu versenden einfach das Tutorial von vorne durchlaufen.

- mac os funktionalität

### Closed:
- versenden von Dateien über tcp
    - windows
    - linux
- GUI Interface zu Auswählen der zu versenden Datei
- Gui Interface zum Auswählen des Ortes zum Speichern der Datei

## Anleitung file Transfer
1. In ./chat0815/fileTransfer/ navigieren
2. ./Client/main.go bearbeiten und gewünschte IP in Zeile 24 hinzufügen
3. ./Server/main.go starten
4. ./Client/main.go starten
Auf dem Bildschirm öffnen sich nun die Fenster um Datei bzw Ort Auszuwählen
5. Zuerst beim "Server" Datei Auswählen und öffnen
6. Danach beim "Client" Zielort auswählen
7. "Client" Fenster schließen
8. "Server" Fenster schließen


EOF