// Заглушка под шифрование пакетов (AES-GCM через WebCrypto).
const te = new TextEncoder();
const td = new TextDecoder();

function b64encode(buf) {
  const bytes = new Uint8Array(buf);
  let bin = '';
  for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
  return btoa(bin);
}
function b64decode(str) {
  const bin = atob(str);
  const bytes = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
  return bytes.buffer;
}

export async function importAesKey(raw32bytes) {
  return crypto.subtle.importKey('raw', raw32bytes, { name: 'AES-GCM' }, false, ['encrypt','decrypt']);
}

export async function encryptPacket(key, obj) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const plaintext = te.encode(JSON.stringify(obj));
  const ciphertext = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, plaintext);
  return { v: 1, iv: b64encode(iv), data: b64encode(ciphertext) };
}

export async function decryptPacket(key, packet) {
  const iv = new Uint8Array(b64decode(packet.iv));
  const data = b64decode(packet.data);
  const plaintext = await crypto.subtle.decrypt({ name: 'AES-GCM', iv }, key, data);
  return JSON.parse(td.decode(plaintext));
}
