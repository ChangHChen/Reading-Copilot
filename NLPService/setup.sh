
if [ ! -d "tls" ]; then
    mkdir "tls"
fi

cd "tls"

if ! dpkg -s libnss3-tools >/dev/null 2>&1; then
    echo "Installing libnss3-tools..."
    sudo apt update
    sudo apt install libnss3-tools -y
fi

if ! command -v mkcert >/dev/null 2>&1; then
    echo "Installing mkcert..."
    git clone https://github.com/FiloSottile/mkcert
    cd mkcert
    go build -ldflags "-X main.Version=$(git describe --tags)"
    sudo cp mkcert /usr/local/bin/
    mkcert -install
    cd ..
    rm -rf mkcert
fi

if [ -f "cert.pem" ] && [ -f "key.pem" ]; then
    echo "Certificate files already exist. No need to regenerate."
    exit 0
fi

mkcert -cert-file cert.pem -key-file key.pem localhost

echo "Certificates generated successfully."