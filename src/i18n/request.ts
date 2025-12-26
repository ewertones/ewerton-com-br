import {getRequestConfig} from "next-intl/server";
import {routing} from "./routing";
import {deepMerge} from "../lib/utils";

export default getRequestConfig(async ({requestLocale}) => {
  let locale = await requestLocale;

  if (!locale || !routing.locales.includes(locale as "en" | "pt-br")) {
    locale = routing.defaultLocale;
  }

  return {
    locale: locale as "en" | "pt-br",
    messages: (await (async () => {
      const userMessages = (await import(`../../messages/${locale}.json`)).default;
      const defaultMessages = (await import(`../../messages/${routing.defaultLocale}.json`)).default;
      
      if (locale === routing.defaultLocale) {
        return defaultMessages;
      }
      
      return deepMerge(defaultMessages, userMessages);
    })())
  };
});
