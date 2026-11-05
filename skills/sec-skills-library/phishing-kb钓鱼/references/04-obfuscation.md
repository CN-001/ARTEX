# Obfuscation & Encryption Reference

## Base64 Obfuscation

### Encode/Decode Basics
```javascript
var encoded = btoa("Hello World");
var decoded = atob(encoded);

var encoded = btoa(encodeURIComponent("Hello, 世界"));
var decoded = decodeURIComponent(atob(encoded));
```

### Base64 Encode Entire HTML Page, Render with document.write
```html
<script>
document.write(decodeURIComponent(atob("BASE64_ENCODED_HTML_HERE")));
</script>
```

### Base64 Encode JS Code, Execute with eval
```html
<head></head>
<body></body>
<script>
eval(decodeURIComponent(atob("BASE64_ENCODED_JS_HERE")));
</script>
```

### Base64 Encode a Function, eval + invoke
```html
<head></head>
<body></body>
<script>
eval(decodeURIComponent(atob("BASE64_ENCODED_FUNCTION_DEFINITION")));
myFunction(); // explicitly invoke after eval
</script>
```

---

## XOR Obfuscation

### Single-Key XOR with URL Encoding
```javascript
function xorEncryptDecrypt(input, key, encode) {
    let output = '';
    const keyChar = key.charCodeAt(0);
    if (!encode) { input = decodeURIComponent(input); }
    for (let i = 0; i < input.length; i++) {
        output += String.fromCharCode(input.charCodeAt(i) ^ keyChar);
    }
    if (encode) { output = encodeURIComponent(output); }
    return output;
}

var enc = xorEncryptDecrypt(htmlContent, "K", true);
var dec = xorEncryptDecrypt(enc, "K", false);
document.write(dec);
```

### Multi-Key XOR
```javascript
function xorEncryptDecrypt(input, key, encode) {
    let output = '';
    let keyIndex = 0;
    if (!encode) { input = decodeURIComponent(input); }
    for (let i = 0; i < input.length; i++) {
        output += String.fromCharCode(input.charCodeAt(i) ^ key.charCodeAt(keyIndex));
        keyIndex = (keyIndex + 1) % key.length;
    }
    if (encode) { output = encodeURIComponent(output); }
    return output;
}
```

### XOR-Encrypted Page Template
```html
<!DOCTYPE html>
<html>
<script>
const enc = `ENCRYPTED_CONTENT_HERE`;
function xorDecrypt(input, key) {
    let output = '';
    const keyChar = key.charCodeAt(0);
    input = decodeURIComponent(input);
    for (let i = 0; i < input.length; i++) {
        output += String.fromCharCode(input.charCodeAt(i) ^ keyChar);
    }
    return output;
}
document.write(xorDecrypt(enc, "E"));
</script>
</html>
```

---

## Dynamic XOR Encryption (Server-Side Key Generation)

The key changes on every page load — no static signature possible.

### PHP: Encrypt + Serve
```php
<?php
$key = bin2hex(random_bytes(16));
$htmlContent = file_get_contents('/var/www/login.html');

function xorEncryptDecrypt($input, $key, $encode = true) {
    $output = '';
    $keyIndex = 0;
    $keyLength = strlen($key);
    if (!$encode) { $input = base64_decode($input); }
    for ($i = 0; $i < strlen($input); $i++) {
        $output .= chr(ord($input[$i]) ^ ord($key[$keyIndex]));
        $keyIndex = ($keyIndex + 1) % $keyLength;
    }
    if ($encode) { $output = base64_encode($output); }
    return $output;
}

$encrypted = xorEncryptDecrypt($htmlContent, $key, true);
?>
<!DOCTYPE html>
<script>
var enc = '<?php echo $encrypted; ?>';
var key = '<?php echo $key; ?>';
function xorDecrypt(input, key) {
    let output = '';
    let keyIndex = 0;
    input = atob(input);
    for (let i = 0; i < input.length; i++) {
        output += String.fromCharCode(input.charCodeAt(i) ^ key.charCodeAt(keyIndex));
        keyIndex = (keyIndex + 1) % key.length;
    }
    return output;
}
document.write(xorDecrypt(enc, key));
</script>
```

---

## AES Encryption (Web Crypto API)

### Encryptor — generates hex key, IV, and ciphertext
```javascript
function bufferToHex(buffer) {
    return Array.from(new Uint8Array(buffer), b => b.toString(16).padStart(2, '0')).join('');
}

async function aesEncrypt(plaintext) {
    let iv = window.crypto.getRandomValues(new Uint8Array(16));
    let key = await window.crypto.subtle.generateKey(
        { name: "AES-CBC", length: 256 }, true, ["encrypt"]
    );
    let encoded = new TextEncoder().encode(plaintext);
    let encrypted = await window.crypto.subtle.encrypt({ name: "AES-CBC", iv: iv }, key, encoded);
    let rawKey = await window.crypto.subtle.exportKey("raw", key);

    console.log("Key:", bufferToHex(rawKey));
    console.log("IV:", bufferToHex(iv));
    console.log("Encrypted:", bufferToHex(encrypted));
}
```

### Decryptor — renders encrypted page
```html
<!DOCTYPE html>
<html><script>
const keyHex = "KEY_HEX_HERE";
const ivHex = "IV_HEX_HERE";
const encHex = "ENCRYPTED_HEX_HERE";

function hexToBuffer(hex) {
    const result = new Uint8Array(hex.length / 2);
    for (let i = 0; i < hex.length; i += 2)
        result[i/2] = parseInt(hex.substring(i, i+2), 16);
    return result.buffer;
}

async function aesDecrypt(keyHex, ivHex, encHex) {
    const key = await window.crypto.subtle.importKey(
        'raw', hexToBuffer(keyHex), { name: 'AES-CBC', length: 256 }, false, ['decrypt']
    );
    const decrypted = await window.crypto.subtle.decrypt(
        { name: 'AES-CBC', iv: hexToBuffer(ivHex) }, key, hexToBuffer(encHex)
    );
    return new TextDecoder().decode(decrypted);
}

aesDecrypt(keyHex, ivHex, encHex).then(text => document.write(text));
</script></html>
```

---

## AES Encryption (aes-js Library)

### Encrypt
```javascript
function aesJsEncrypt(plaintext) {
    var bytes = aesjs.utils.utf8.toBytes(plaintext);
    var padded = aesjs.padding.pkcs7.pad(bytes);
    var iv = [21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36];
    var key = [0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31];
    var aesCbc = new aesjs.ModeOfOperation.cbc(key, iv);
    var encrypted = aesCbc.encrypt(padded);
    console.log(aesjs.utils.hex.fromBytes(key));
    console.log(aesjs.utils.hex.fromBytes(iv));
    console.log(aesjs.utils.hex.fromBytes(encrypted));
}
```

### Decrypt + Render
```javascript
function aesJsDecrypt(keyHex, ivHex, encHex) {
    var keyBytes = aesjs.utils.hex.toBytes(keyHex);
    var ivBytes = aesjs.utils.hex.toBytes(ivHex);
    var encBytes = aesjs.utils.hex.toBytes(encHex);
    var aesCbc = new aesjs.ModeOfOperation.cbc(keyBytes, ivBytes);
    var decrypted = aesCbc.decrypt(encBytes);
    var unpadded = aesjs.padding.pkcs7.strip(decrypted);
    document.write(aesjs.utils.utf8.fromBytes(unpadded));
}
```

Include the library:
```html
<script src="https://cdnjs.cloudflare.com/ajax/libs/aes-js/3.1.2/index.min.js"></script>
```
