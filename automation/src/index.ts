import { config } from "dotenv";
import { Growi } from "./growi";

async function main() {
  // load .env config
  config();

  const growi = await Growi.createInstance();

  try {
    await growi.addDummyPage(100);
  } catch (err) {
    console.log(err);
  } finally {
    await growi.close();
  }
}

main().catch((err) => {
  console.log(err);
});
