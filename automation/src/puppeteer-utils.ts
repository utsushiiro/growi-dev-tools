import puppeteer, { Page, FrameBase } from "puppeteer";

export type CustomPage = {
  getPropertyValueBySelector: (selector: string, property: string) => Promise<string | undefined>;
  replaceInputValue: (selector: string, value: string) => Promise<unknown>;
  replaceTextAreaValue: (selector: string, value: string) => Promise<unknown>;
} & Page;

export async function openPageInNewBrowser(isHeadless: boolean = false) {
  const width = 1200;
  const height = 1000;

  const browser = await puppeteer.launch({
    ignoreHTTPSErrors: true,
    headless: isHeadless,
    // ref: https://github.com/puppeteer/puppeteer/issues/1183#issuecomment-363569401
    args: [`--window-size=${width},${height}`],
  });

  const page = (await browser.pages())[0];
  await page.setViewport({ width, height });

  const customPage: CustomPage = page as any;
  customPage.getPropertyValueBySelector = getPropertyValueBySelector.bind(null, page);
  customPage.replaceInputValue = replaceInputValue.bind(null, page);
  customPage.replaceTextAreaValue = replaceTextAreaValue.bind(null, page);

  return { browser, page: customPage };
}

export async function getPropertyValueBySelector(frameBase: FrameBase, selector: string, property: string) {
  return (await (await frameBase.$(selector))?.getProperty(property))?.jsonValue() as Promise<string | undefined>;
}

export async function replaceInputValue(page: Page, selector: string, value: string) {
  await page.evaluate(
    ({ selector, value }) => {
      const element = document.querySelector(selector);
      const nativeInputValueSetter = (Object.getOwnPropertyDescriptor(
        window.HTMLInputElement.prototype,
        "value"
      ) as any).set;
      nativeInputValueSetter.call(element, value);
      element.dispatchEvent(new Event("input", { bubbles: true }));
    },
    { selector, value }
  );
}

export async function replaceTextAreaValue(page: Page, selector: string, value: string) {
  await page.evaluate(
    ({ selector, value }) => {
      const element = document.querySelector(selector);
      const nativeTextAreaValueSetter = (Object.getOwnPropertyDescriptor(
        window.HTMLTextAreaElement.prototype,
        "value"
      ) as any).set;
      nativeTextAreaValueSetter.call(element, value);
      element.dispatchEvent(new Event("input", { bubbles: true }));
    },
    { selector, value }
  );
}
