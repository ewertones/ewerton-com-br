import createMiddleware from "next-intl/middleware";
import {routing} from "./i18n/routing";
import {NextRequest, NextResponse} from "next/server";
import {LINKEDIN_URL, RESUME_PATH} from "./constants";

export default async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;

  if (pathname === '/linkedin') {
    return NextResponse.redirect(LINKEDIN_URL);
  }

  if (pathname === '/cv') {
    return NextResponse.redirect(new URL(RESUME_PATH, request.url));
  }

  const handleI18nRouting = createMiddleware(routing);
  return handleI18nRouting(request);
}

export const config = {
  matcher: [
    "/((?!api|_next|_vercel|opengraph-image|twitter-image|favicon|icon|apple-icon|apple-touch-icon|manifest|robots|sitemap|.*\\..*).*)"
  ]
};
