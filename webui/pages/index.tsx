import Head from "next/head"
import { useRouter } from "next/router"

import { Layout } from "@/components/layout"
import Cookies from 'js-cookie';


export default function IndexPage() {
  const router = useRouter()
  const token = router.query.token
  if (token) {
    localStorage.setItem("token", token as string)
    Cookies.set('token', token as string);
    router.push("/")
  }

  return (
    <Layout>
      <Head>
        <title>NotionBoy</title>
        <meta
          name="description"
          content="NotionBoy is a Note-taking app built with Notion API."
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <section className="container grid items-center gap-6 pt-6 pb-8 md:py-10">
        <div className="">
          <h1>Welcome to NotionBoy</h1>
        </div>
      </section>
    </Layout>
  )
}
