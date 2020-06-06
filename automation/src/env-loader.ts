export default function loadEnv(key: string) {
  if (process.env[key] == null) {
    throw new Error(`環境変数(${key})に値が設定されていません`);
  }
  return process.env[key] as string;
}
