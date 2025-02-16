# GomniRun - Cross-Platform Script Runner 🏃‍♂️💨

GomniRun is a **cross-platform script execution tool** with a **GUI (Fyne) and CLI**.  
It allows you to set **variables, commands, and arguments dynamically** and saves them in a `config.json` file.

## 📌 Features

- ✅ **Cross-platform GUI** (Fyne) & **CLI** support.
- ✅ **Stores last-used variables** in `config.json`.
- ✅ **Supports different variable types** (`string`, `bool`, `file`, `directory`).
- ✅ **Fully configurable command execution**.
- ✅ **Dark & Light mode support** with high-contrast output.
- ✅ **GitHub integration & open-source!**

## 📸 Screenshots

*(Add screenshots here if possible!)*

---

## ⚡ Installation

### **1️⃣ Install Go & Dependencies**

Ensure you have Go **installed (Go 1.21+)**:

```sh
go version
```

If not installed, download it from [Go's official site](https://go.dev/dl/).

Then, install the **Fyne** dependency:

```sh
go get fyne.io/fyne/v2
```

### **2️⃣ Clone the Repo**

```sh
git clone https://github.com/<your-username>/GomniRun.git
cd GomniRun
```

### **3️⃣ Run the App**

#### GUI Mode

```sh
go run cmd/fyne-ui/main.go
```

#### CLI Mode

```sh
go run cmd/cli/main.go
```

---

## ⚙️ Configuration (`config.json`)

GomniRun reads from a **config.json** file that stores **command templates and variables**.

Example `config.json`:

```json
{
  "command": "bash {script} -variable={var1} -other-variable={var2}",
  "variables": {
    "script": { "type": "file", "value": "./test_script.sh" },
    "var1": { "type": "string", "value": "hello" },
    "var2": { "type": "string", "value": "world" }
  },
  "overwrite": true
}
```

### 📝 **How It Works**

- `{script}` → Points to the script file.
- `{var1}` → Holds the first variable (`hello`).
- `{var2}` → Holds the second variable (`world`).
- `"overwrite": true` → Saves last-used values automatically.

---

## 🏃‍♂️ Running a Test Script

You can create a simple test script (`test_script.sh`) to check execution:

```sh
#!/bin/bash
echo "🚀 Test Script Executed!"
echo "Variable: $1"
echo "Other Variable: $2"
```

Make it executable:

```sh
chmod +x test_script.sh
```

Update `config.json` and run the app!

---

## Built With

```sh
go build -o builds/gomnirun main.go
go build -o builds/gomnirun-cli ./cmd/cli/
go build -o builds/gomnirun-ui ./cmd/fyne-ui/

```


## 🌎 Contributing

🔥 **Want to improve GomniRun?** Contributions are welcome!  

- Fork the repo  
- Create a feature branch  
- Submit a PR! 🎉  

---

## 📜 License

This project is licensed under the **MIT License**.  
Feel free to use and modify it as needed!

---

## ⭐ Support the Project

If you like this project, give it a ⭐ on GitHub!  
Happy coding! 🚀🔥
