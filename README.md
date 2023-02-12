<p align="right">
بِسْــــــــــــــمِ اللَّهِ الرَّحْمَنِ الرَّحِيم 
</p>


# WA Bot dengan Asisten ChatGPT
Projek sederhana membuat bot Whatsapp dengan asisten AI ChatGPT.

## How to use
- clone atau download repo ini
- buat file .env
    ```bash
    cp .env.example .env
    ```
- Copy API keys openAI https://beta.openai.com/account/api-keys paste di env pada API_KEY.
- Build server
    ```bash
    go build .
    ```
- Jalankan server
    ```bash
    ./wa-bot-chatgpt
    ```
- Pastikan waktu mengirim pesan menggunakan trigger **-ask**
    ```
    -ask buatkan kalimat perkenalan
    ```