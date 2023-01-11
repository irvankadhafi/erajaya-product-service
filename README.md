# Erajaya Product Service API

## Design Architecture
Konsep yang digunakan pada sevice ini ialah clean architecture. 

Tiap komponen tidak bergantung pada framework ataupun database yang digunakan (independen). Konsep ini di kemukakan oleh Uncle Bob yang bisa dibaca pada artikel berikut https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html.



```bash
.
├── cache/
|   # package yang berguna untuk memanage cache, seperti GET, STORE, dll
|   
├── db/
|   # berisikan skrip migrasi database dalam bentul file sql
|   
├── internal/
|   # Berisi file private aplikasi dan library.
│   ├── config/
│   │   # Menyimpan file konfigurasi dan default value
│   ├── console/
│   ├── db/
│   └── delivery/
│   │   # Layer ini bertugas sebagai presenter atau menyajikan output ke client
│   │   # Banyak metode yang dapat digunakan, seperti: HTTP REST API, gRPC, GraphQL. Pada kasus ini saya menggunakan HTTP REST API
│   └── helper/
│   └── model/
│   │   # Layer ini menyimpan model yang akan digunakan pada layer lainnya. 
│   │   # Layer ini dapat diakses oleh semua layer.
│   └── repository/
│   │   # Layer ini menyimpan database handler dan cache handler. 
│   │   # Tidak ada business logic di layer ini.
│   │   # Bertugas untuk menentukan datastore apa yang digunakan, pada kasus ini saya menggunakan RDBMS PostgreSQL
│   └── usecase/
│       # Layer ini berisi business logic pada domain.
│       # Melakukan kontrol, sehingga bisa memilih repository yang akan digunakan
│       # Bertugas sebagai penghubung antara layer repository dengan layer delivery.
│       # Proses validasi terjadi pada layer ini
|   
├── utils/
|   # 
├── config.yml
|   # config untuk menjalankan server
├── go.mod
|   # file go.mod dipergunakan oleh go module (jika go mod diaktifkan).
├── main.go
|   # Program main utama
├── Makefile
|   # file Makefile dipergunakan oleh command `make`.
└── ...
```

## List API
Berikut fitur dan endpoint API yang terdapat dalam project ini :

| Method | Feature        | Endpoint   |
|--------|----------------|------------|
| POST   | Create product | /products/ |
| GET    | Search product | /products/ |
