echo 'Installer made with GISC'
if ! [ -x "$(command -v go)" ]; then
	 apt install golang-go>&2
	exit 1
fi
if ! [ -x "$(command -v git)" ]; then
	 apt install git>&2
	exit 1
fi

go build || exit 1
