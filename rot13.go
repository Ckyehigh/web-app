<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Encrypt/Decrypt</title>
    <style>
        body {
            background-color: black;
            color: darkred;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
            font-family: 'Courier New', Courier, monospace;
        }
        .container {
            text-align: center;
            padding: 20px;
            background-color: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
        }
        label, input, select {
            display: block;
            margin: 10px auto;
            width: 80%;
            padding: 10px;
            border: 1px solid darkred;
            border-radius: 5px;
            background-color: #333; /* Solid dark background for input fields */
            color: darkred;
        }
        input[type="submit"] {
            background-color: darkred;
            color: white;
            cursor: pointer;
            width: 60%;
        }
        input[type="submit"]:hover {
            background-color: #a52a2a;
        }
        h1 {
            margin-bottom: 20px;
        }
        h2 {
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Encrypt/Decrypt Text</h1>
        <form action="/process" method="post">
            <label for="text">Text:</label>
            <input type="text" id="text" name="text" value="{{.Text}}" required>
            <label for="key">Key:</label>
            <input type="text" id="key" name="key" value="{{.Key}}" required>
            <label for="action">Action:</label>
            <select id="action" name="action" required>
                <option value="E" {{if eq .Action "E"}}selected{{end}}>Encrypt</option>
                <option value="D" {{if eq .Action "D"}}selected{{end}}>Decrypt</option>
            </select>
            <input type="submit" value="Submit">
        </form>

        {{if .Result}}
        <h2>Result</h2>
        <p>{{.Result}}</p>
        {{end}}
    </div>
</body>
</html>
