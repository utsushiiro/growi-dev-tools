import { Browser } from "puppeteer";
import loadEnv from "./env-loader";
import { openPageInNewBrowser, replaceInputValue, CustomPage } from "./puppeteer-utils";
import { LoremIpsum } from "lorem-ipsum";
import { v4 as uuid } from "uuid";

export class Growi {
  private dummyTextGenerator: LoremIpsum;

  private constructor(private browser: Browser, private page: CustomPage) {
    this.dummyTextGenerator = new LoremIpsum({
      sentencesPerParagraph: {
        max: 30,
        min: 5,
      },
      wordsPerSentence: {
        max: 30,
        min: 5,
      },
    });
  }

  static async createInstance(): Promise<Growi> {
    const adminUserName = loadEnv("GROWI_ADMIN_USER_NAME");
    const adminUserPassword = loadEnv("GROWI_ADMIN_USER_PASSWORD");
    const growiURL = loadEnv("GROWI_URL");

    const { browser, page } = await openPageInNewBrowser(false);
    await page.goto(growiURL);

    // ログイン
    await page.waitForSelector("#login-form");
    await replaceInputValue(page, "input[name='loginForm[username]']", adminUserName);
    await replaceInputValue(page, "input[name='loginForm[password]']", adminUserPassword);
    await Promise.all([page.waitForNavigation({ waitUntil: ["load", "networkidle2"] }), page.click("#login")]);

    return new Growi(browser, page);
  }

  async close() {
    await this.browser.close();
  }

  /**
   * /dummy にダミーのページを指定した数だけ作成する
   * @param num ダミーページ数
   */
  async addDummyPage(num: number) {
    for (let i = 0; i < num; i++) {
      await this.page.click("#create-page-button");

      // モーダルが表示されるまで待つ
      await this.page.waitForSelector(".grw-create-page .modal-body .row:nth-child(2) input", { visible: true });

      // 適当なページ名を設定
      const pageName = `/dummy/${uuid()}`;
      await this.page.replaceInputValue(".grw-create-page .modal-body .row:nth-child(2) input", pageName);

      // Createボタンをクリック -> ページ作成画面に遷移
      await Promise.all([
        this.page.waitForNavigation({ waitUntil: ["load", "networkidle2"] }),
        this.page.click(".grw-create-page .modal-body .row:nth-child(2) button"),
      ]);

      // CodeMirrorのエディタが表示されるまで待つ
      await this.page.waitForSelector(".CodeMirror", { visible: true });

      // 内容を書き換える
      // https://stackoverflow.com/questions/21844574/programmatically-edit-codemirror-contents-without-access-to-object
      const dummyText = this.dummyTextGenerator.generateParagraphs(10);
      await this.page.evaluate((dummyText) => {
        const codeMirrorDiv = document.querySelector(".CodeMirror");
        (codeMirrorDiv as any).CodeMirror.setValue(dummyText);
      }, dummyText);

      // Createボタンをクリック -> 作成したページに遷移
      await Promise.all([
        this.page.waitForNavigation({ waitUntil: ["load", "networkidle2"] }),
        this.page.click("#caret"),
      ]);

      console.log(`[${i + 1}/${num}] page created: ${pageName}`);
    }
  }
}
