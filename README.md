# Zip Downloader

Build like so

```bash
go build -o zipdl
```

to install this command add it to your PATH or copy it like so `sudo cp zipdl /usr/local/bin/`

The program takes three arguments

 - `url` the url we want to download the zip from
 - `path` the path where we want to store the unzipped files
 - `interval` (optional) this causes the program to run forever and specifies how often we want to download and extract the files eg. `15m`, `6h`, `1d`, `1h30m`

I use this to download pricing data for Elder Scrolls Online for an addon called Tamriel Trade Center. I play on Linux so there are not a lot of solutions out there. Here is how I use it:
```bash
zipdl -path="/home/uberswe/.steam/steam/steamapps/compatdata/306130/pfx/drive_c/users/steamuser/My Documents/Elder Scrolls Online/live/AddOns/TamrielTradeCentre" -interval="6h" -url="https://eu.tamrieltradecentre.com/download/PriceTable"
```