import { NavItem } from "@/types/nav"

interface SiteConfig {
  name: string
  description: string
  mainNav: NavItem[]
  links: {
    twitter: string
    github: string
    docs: string
  }
}

export const siteConfig: SiteConfig = {
  name: "NotionBoy",
  description:
    "Beautifully designed components built with Radix UI and Tailwind CSS.",
  mainNav: [
    {
      title: "Home",
      href: "/web",
    },
    {
      title: "Chat",
      href: "/web/chat.html",
    },
  ],
  links: {
    twitter: "https://twitter.com/LiuVaayne",
    github: "https://github.com/Vaayne/NotionBoy",
    docs: "https://www.theboys.tech/",
  },
}
