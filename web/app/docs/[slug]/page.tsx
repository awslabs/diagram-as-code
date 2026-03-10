import DocsArticlePage from '@/components/docs/DocsArticlePage'
import { docsSectionIds, type DocsSectionId } from '@/lib/docs'
import { notFound } from 'next/navigation'

export function generateStaticParams() {
  return docsSectionIds.map((slug) => ({ slug }))
}

export default async function DocsArticleRoute({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = await params

  if (!docsSectionIds.includes(slug as DocsSectionId)) {
    notFound()
  }

  return <DocsArticlePage slug={slug as DocsSectionId} />
}
