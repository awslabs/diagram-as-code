'use client'

import { createContext, useContext, useState, useEffect, ReactNode } from 'react'

export type Lang = 'en' | 'pt'

// ── English ───────────────────────────────────────────────────────────────────

const en = {
  // Editor page
  loadingEditor: 'Loading editor…',
  examples: 'Examples',
  generate: 'Generate',
  docs: 'Docs',
  yamlEditor: 'YAML Editor',
  lines: (n: number) => `${n} lines`,
  preview: 'Preview',
  mode: (fmt: string) => `${fmt.toUpperCase()} mode · Ctrl+Enter to generate`,
  footerBy: 'by Fernando Azevedo',

  // DiagramPreview
  generating: 'Generating diagram…',
  generationFailed: 'Generation failed',
  downloadPng: 'Download PNG',
  downloadPdf: 'Download PDF',
  drawioReady: 'draw.io file ready',
  drawioHelper: 'Open with draw.io or diagrams.net',
  drawioEmbedded: 'Editing draw.io directly inside the app',
  downloadDrawio: 'Download .drawio',
  pdfReady: 'PDF file ready',
  pdfHelper: 'Download the generated PDF document',
  emptyStatePre: 'Click',
  emptyStatePost: 'to preview your diagram',
  generatedDiagram: 'Generated diagram',

  // Docs — navigation
  backToEditor: 'Back to Editor',
  documentation: 'Documentation',
  contents: 'Contents',
  navOverview: 'Overview',
  navWebEditor: 'Web Editor',
  navCli: 'CLI Usage',
  navDrawio: 'Draw.io Export',
  navApi: 'API Reference',
  navMcp: 'MCP Server',
  navLocalDev: 'Local Development',
  navExamples: 'YAML Examples',
  navAbout: 'About the Author',
  viewOnGithub: 'View on GitHub',

  // Docs — Overview
  overviewP1: 'diagram-as-code is an open-source CLI tool that generates AWS architecture diagrams from YAML. Write human-readable YAML, get pixel-perfect PNG diagrams or draw.io files — all without touching any GUI diagramming tool.',
  overviewP2Fork: 'This fork, maintained by',
  overviewP2Extends: ', extends the original',
  overviewP2With: 'with a browser-based editor hosted on Vercel and a native draw.io export pipeline.',
  overviewForkRepo: 'Fork repository:',
  whatToolProduces: 'What this tool produces',
  outputCol: 'Output',
  formatCol: 'Format',
  useCaseCol: 'Use case',
  pngDiagram: 'PNG diagram',
  pngUseCase: 'Docs, wikis, PRs, CI artifacts',
  drawioFile: 'draw.io file',
  drawioUseCase: 'Editable diagrams in diagrams.net',

  // Docs — Web Editor
  webEditorP1: 'The web editor lets you write YAML and generate diagrams directly in the browser. No installation needed. Powered by Vercel serverless functions running the same Go engine as the CLI.',
  howToUse: 'How to use it',
  webEditorStep1: 'Open the editor at',
  webEditorStep2: 'Write or paste your YAML in the left panel, or pick an example from the',
  webEditorStep2b: 'dropdown',
  webEditorStep3: 'Choose output format:',
  webEditorStep4: 'Click',
  webEditorStep4b: 'or press',
  webEditorStep5: 'Download the result with the button in the preview panel',
  webEditorCallout: 'The editor uses Monaco (the VS Code engine) with YAML syntax highlighting. draw.io files can be opened in',
  webEditorCalloutSuffix: 'for further editing.',

  // Docs — CLI
  install: 'Install',
  basicUsage: 'Basic usage',
  allFlags: 'All flags',
  flagCol: 'Flag',
  defaultCol: 'Default',
  descriptionCol: 'Description',
  flagOutput: 'Output file name',
  flagDrawio: 'Generate draw.io file instead of PNG',
  flagForce: 'Overwrite output without confirmation',
  flagVerbose: 'Enable verbose logging',
  flagWidth: 'Resize output image width (PNG only)',
  flagHeight: 'Resize output image height (PNG only)',
  flagTemplate: 'Process input as Go text/template',
  flagCfn: '[Beta] Create diagram from CloudFormation template',
  flagUntrusted: 'Allow definition files from non-official URLs',

  // Docs — Draw.io
  drawioP1: 'The draw.io export uses the same resource/link model and layout engine as PNG rendering. Resources become',
  drawioP1b: 'nodes; links become edges. Official AWS SVG icons are embedded as data URIs so the file is fully self-contained.',
  exportPipeline: 'Export pipeline',
  drawioStep1: 'YAML is parsed into the resource graph',
  drawioStep2: 'The same layout pass runs',
  drawioStep3: 'Children are reordered by link topology (matches PNG ordering)',
  drawioStep4: 'Leaf resources get AWS icons embedded as base64 data URIs',
  drawioStep5: 'Groups get AWS4 group styles with correct fill/stroke colors',
  drawioStep6: 'Links become',
  drawioStep6b: 'edges with optional source/target labels',
  drawioStep7: 'Output is written as',
  drawioStep7b: 'XML',
  drawioCallout: 'Open the generated',
  drawioCalloutAt: 'file at',
  drawioCalloutSuffix: '— all icons, labels and group borders render without internet access.',

  // Docs — API
  apiP1: 'The serverless function at',
  apiP1b: 'powers the web editor and can be called directly from any HTTP client or AI agent.',
  requestBody: 'Request body',
  response: 'Response',
  curlExample: 'curl example',
  jsExample: 'JavaScript / TypeScript example',
  propertyCol: 'Property',
  valueCol: 'Value',
  bodyCol: 'Body',

  // Docs — MCP
  mcpP1: 'The MCP (Model Context Protocol) server lets AI assistants like Claude generate AWS architecture diagrams directly from chat — no CLI required.',
  installMcp: 'Install the MCP server',
  configureClaude: 'Configure in Claude Desktop',
  configureClaudeP: 'Add the following to your',
  configureClaudeP2: 'On macOS the config file is at',
  availableMcpTools: 'Available MCP tools',
  toolCol: 'Tool',
  mcpTool1Desc: 'Generate a PNG diagram from a DAC YAML string',
  mcpTool2Desc: 'Generate a draw.io XML file from a DAC YAML string',
  examplePrompt: 'Example prompt for Claude',
  mcpCallout: 'Claude will use the MCP server to invoke the Go diagram engine, applying the same layout and icon rules as the CLI tool.',

  // Docs — Local Dev
  prerequisites: 'Prerequisites',
  versionCol: 'Version',
  cloneAndRun: 'Clone and run the full stack locally',
  localDevP: 'Open',
  localDevPSuffix: '— the editor connects to your local Go server automatically. No Vercel CLI needed.',
  runTests: 'Run tests',
  buildCli: 'Build the CLI',
  runWithVercel: 'Run with Vercel CLI (full stack)',

  // Docs — Examples
  minimalS3: 'Minimal — S3 bucket',
  albEc2: 'ALB + EC2 in a VPC',
  keyYamlConcepts: 'Key YAML concepts',
  fieldCol: 'Field',
  yamlConcept1: 'URLs or local files defining AWS resource types and icons',
  yamlConcept2: 'AWS resource type (e.g. AWS::EC2::Instance) or diagram element',
  yamlConcept3: 'List of resource names placed inside this container',
  yamlConcept4: '"vertical" or "horizontal" layout for children',
  yamlConcept5: 'Named preset from the definition file (e.g. "PublicSubnet")',
  yamlConcept6: 'Override the display label',
  yamlConcept7: 'Place resources on the border of a container',
  yamlConcept8: 'Resource names to connect',
  yamlConcept9: 'Label on the source side of the edge',
  yamlConcept10: '"Open" = open arrow, "none" = no arrow',

  // Docs — About
  authorSubtitle: 'Senior Solutions Architect · 16+ years global experience',
  authorBio1: 'Senior Solutions Architect with 16+ years of global experience delivering impactful, secure, and scalable digital solutions. Specialized in designing scalable, secure, and cost-efficient cloud architectures. Currently at',
  authorBio1b: ', bridging business goals and technology through innovation and best practices.',
  authorBio2: 'Deep expertise in',
  authorBio2b: ', Event-Driven systems, Microservices patterns, and high-scale architectures. Experienced with Data Mesh, Kafka, EKS, and AWS Well-Architected Framework.',
  authorBio3: 'Security-first development approach with extensive experience in CI/CD, Infrastructure as Code (Terraform, CDK), GitOps, and Zero Trust Architecture. Compliance with PCI-DSS, ISO 27001, and GDPR.',
  footerFork: 'diagram-as-code · Fork by Fernando Azevedo',

  // Builder page
  builderLink: 'Builder',
  builderTitle: 'YAML Builder',
  builderDesc: 'Build your AWS diagram visually — no YAML required',
  defFileSection: 'Definition File',
  officialUrl: 'Official URL (recommended)',
  customUrlLabel: 'Custom URL',
  localFileLabel: 'Local File',
  filePathLabel: 'File path',
  canvasSection: 'Canvas Settings',
  canvasDir: 'Layout Direction',
  resourcesSection: 'Resources',
  addResource: '+ Add Resource',
  resourceName: 'Name (YAML key)',
  resourceType: 'Type',
  resourceTitle: 'Title (optional)',
  resourcePreset: 'Preset (optional)',
  resourceDirection: 'Direction (optional)',
  resourceAlign: 'Align (optional)',
  resourceChildren: 'Children',
  borderChildrenLabel: 'Border Children',
  addBorderChild: '+ Add border child',
  positionLabel: 'Position',
  removeLabel: 'Remove',
  linksSection: 'Links',
  addLink: '+ Add Link',
  linkSource: 'Source',
  linkTarget: 'Target',
  sourcePosLabel: 'Source Position (optional)',
  targetPosLabel: 'Target Position (optional)',
  arrowHeadLabel: 'Arrow Head',
  edgeLabels: 'Edge Labels (optional)',
  yamlPreview: 'YAML Preview',
  copyYaml: 'Copy YAML',
  useInEditor: 'Use in Editor',
  copied: 'Copied!',
  noResources: 'No resources yet. Click "+ Add Resource" to start.',
  noLinks: 'No links yet.',
  searchTypes: 'Search types…',
  selectOptional: '(none)',
  openArrow: 'Open',
  noArrow: 'None',
  defaultArrow: 'Default',
  iconFillLabel: 'Icon Fill (optional)',
  resourceNamePlaceholder: 'e.g. MyVPC',
  customUrlPlaceholder: 'https://...',
  localFilePlaceholder: '~/my-definition.yaml',
  linesCount: (n: number) => `${n} lines`,
}

// ── Portuguese ────────────────────────────────────────────────────────────────

const pt: typeof en = {
  loadingEditor: 'Carregando editor…',
  examples: 'Exemplos',
  generate: 'Gerar',
  docs: 'Docs',
  yamlEditor: 'Editor YAML',
  lines: (n: number) => `${n} linhas`,
  preview: 'Visualização',
  mode: (fmt: string) => `Modo ${fmt.toUpperCase()} · Ctrl+Enter para gerar`,
  footerBy: 'por Fernando Azevedo',

  generating: 'Gerando diagrama…',
  generationFailed: 'Falha na geração',
  downloadPng: 'Baixar PNG',
  downloadPdf: 'Baixar PDF',
  drawioReady: 'Arquivo draw.io pronto',
  drawioHelper: 'Abrir com draw.io ou diagrams.net',
  drawioEmbedded: 'Editando o draw.io direto dentro da aplicacao',
  downloadDrawio: 'Baixar .drawio',
  pdfReady: 'Arquivo PDF pronto',
  pdfHelper: 'Baixe o documento PDF gerado',
  emptyStatePre: 'Clique em',
  emptyStatePost: 'para visualizar seu diagrama',
  generatedDiagram: 'Diagrama gerado',

  backToEditor: 'Voltar ao Editor',
  documentation: 'Documentação',
  contents: 'Conteúdo',
  navOverview: 'Visão Geral',
  navWebEditor: 'Editor Web',
  navCli: 'Uso do CLI',
  navDrawio: 'Exportar Draw.io',
  navApi: 'Referência da API',
  navMcp: 'Servidor MCP',
  navLocalDev: 'Desenvolvimento Local',
  navExamples: 'Exemplos YAML',
  navAbout: 'Sobre o Autor',
  viewOnGithub: 'Ver no GitHub',

  overviewP1: 'diagram-as-code é uma ferramenta CLI de código aberto que gera diagramas de arquitetura AWS a partir de YAML. Escreva YAML legível e obtenha diagramas PNG pixel-perfect ou arquivos draw.io — tudo sem precisar de ferramentas GUI de diagramação.',
  overviewP2Fork: 'Este fork, mantido por',
  overviewP2Extends: ', estende o original',
  overviewP2With: 'com um editor baseado em navegador hospedado na Vercel e um pipeline nativo de exportação para draw.io.',
  overviewForkRepo: 'Repositório do fork:',
  whatToolProduces: 'O que esta ferramenta gera',
  outputCol: 'Saída',
  formatCol: 'Formato',
  useCaseCol: 'Caso de uso',
  pngDiagram: 'Diagrama PNG',
  pngUseCase: 'Docs, wikis, PRs, artefatos de CI',
  drawioFile: 'Arquivo draw.io',
  drawioUseCase: 'Diagramas editáveis no diagrams.net',

  webEditorP1: 'O editor web permite escrever YAML e gerar diagramas diretamente no navegador. Sem instalação. Alimentado por funções serverless Vercel rodando o mesmo motor Go do CLI.',
  howToUse: 'Como usar',
  webEditorStep1: 'Abra o editor em',
  webEditorStep2: 'Escreva ou cole seu YAML no painel esquerdo, ou escolha um exemplo no menu',
  webEditorStep2b: '',
  webEditorStep3: 'Escolha o formato de saída:',
  webEditorStep4: 'Clique em',
  webEditorStep4b: 'ou pressione',
  webEditorStep5: 'Baixe o resultado com o botão no painel de visualização',
  webEditorCallout: 'O editor usa Monaco (o motor do VS Code) com destaque de sintaxe YAML. Arquivos draw.io podem ser abertos em',
  webEditorCalloutSuffix: 'para edição posterior.',

  install: 'Instalar',
  basicUsage: 'Uso básico',
  allFlags: 'Todas as flags',
  flagCol: 'Flag',
  defaultCol: 'Padrão',
  descriptionCol: 'Descrição',
  flagOutput: 'Nome do arquivo de saída',
  flagDrawio: 'Gerar arquivo draw.io em vez de PNG',
  flagForce: 'Sobrescrever saída sem confirmação',
  flagVerbose: 'Ativar log detalhado',
  flagWidth: 'Redimensionar largura da imagem (somente PNG)',
  flagHeight: 'Redimensionar altura da imagem (somente PNG)',
  flagTemplate: 'Processar entrada como Go text/template',
  flagCfn: '[Beta] Criar diagrama a partir de template CloudFormation',
  flagUntrusted: 'Permitir arquivos de definição de URLs não-oficiais',

  drawioP1: 'A exportação draw.io usa o mesmo modelo de recurso/link e motor de layout do PNG. Recursos se tornam nós',
  drawioP1b: '; links viram arestas. Ícones SVG oficiais da AWS são embutidos como data URIs para que o arquivo seja totalmente autocontido.',
  exportPipeline: 'Pipeline de exportação',
  drawioStep1: 'O YAML é parseado no grafo de recursos',
  drawioStep2: 'O mesmo passo de layout é executado',
  drawioStep3: 'Os filhos são reordenados por topologia de links (igual ao PNG)',
  drawioStep4: 'Recursos folha recebem ícones AWS embutidos como data URIs base64',
  drawioStep5: 'Grupos recebem estilos AWS4 com cores corretas de preenchimento/borda',
  drawioStep6: 'Links se tornam arestas',
  drawioStep6b: 'com rótulos opcionais de origem/destino',
  drawioStep7: 'A saída é escrita como XML',
  drawioStep7b: '',
  drawioCallout: 'Abra o arquivo',
  drawioCalloutAt: 'gerado em',
  drawioCalloutSuffix: '— todos os ícones, rótulos e bordas de grupos renderizam sem acesso à internet.',

  apiP1: 'A função serverless em',
  apiP1b: 'alimenta o editor web e pode ser chamada diretamente de qualquer cliente HTTP ou agente de IA.',
  requestBody: 'Corpo da requisição',
  response: 'Resposta',
  curlExample: 'Exemplo com curl',
  jsExample: 'Exemplo em JavaScript / TypeScript',
  propertyCol: 'Propriedade',
  valueCol: 'Valor',
  bodyCol: 'Corpo',

  mcpP1: 'O servidor MCP (Model Context Protocol) permite que assistentes de IA como o Claude gerem diagramas de arquitetura AWS diretamente pelo chat — sem precisar do CLI.',
  installMcp: 'Instalar o servidor MCP',
  configureClaude: 'Configurar no Claude Desktop',
  configureClaudeP: 'Adicione o seguinte ao seu',
  configureClaudeP2: 'No macOS, o arquivo de configuração está em',
  availableMcpTools: 'Ferramentas MCP disponíveis',
  toolCol: 'Ferramenta',
  mcpTool1Desc: 'Gera um diagrama PNG a partir de uma string YAML DAC',
  mcpTool2Desc: 'Gera um arquivo draw.io XML a partir de uma string YAML DAC',
  examplePrompt: 'Exemplo de prompt para o Claude',
  mcpCallout: 'O Claude usará o servidor MCP para invocar o motor de diagramas Go, aplicando as mesmas regras de layout e ícones do CLI.',

  prerequisites: 'Pré-requisitos',
  versionCol: 'Versão',
  cloneAndRun: 'Clonar e rodar o stack completo localmente',
  localDevP: 'Abra',
  localDevPSuffix: '— o editor conecta automaticamente ao seu servidor Go local. Sem necessidade do Vercel CLI.',
  runTests: 'Executar testes',
  buildCli: 'Compilar o CLI',
  runWithVercel: 'Rodar com Vercel CLI (stack completo)',

  minimalS3: 'Mínimo — bucket S3',
  albEc2: 'ALB + EC2 em uma VPC',
  keyYamlConcepts: 'Conceitos-chave do YAML',
  fieldCol: 'Campo',
  yamlConcept1: 'URLs ou arquivos locais que definem tipos de recursos AWS e ícones',
  yamlConcept2: 'Tipo de recurso AWS (ex: AWS::EC2::Instance) ou elemento de diagrama',
  yamlConcept3: 'Lista de nomes de recursos colocados dentro deste container',
  yamlConcept4: 'Layout "vertical" ou "horizontal" para os filhos',
  yamlConcept5: 'Preset nomeado do arquivo de definição (ex: "PublicSubnet")',
  yamlConcept6: 'Sobrescreve o rótulo de exibição',
  yamlConcept7: 'Coloca recursos na borda de um container',
  yamlConcept8: 'Nomes dos recursos a conectar',
  yamlConcept9: 'Rótulo no lado de origem da aresta',
  yamlConcept10: '"Open" = seta aberta, "none" = sem seta',

  authorSubtitle: 'Arquiteto de Soluções Sênior · 16+ anos de experiência global',
  authorBio1: 'Arquiteto de Soluções Sênior com mais de 16 anos de experiência global entregando soluções digitais impactantes, seguras e escaláveis. Especializado em projetar arquiteturas cloud escaláveis, seguras e com custo otimizado. Atualmente no',
  authorBio1b: ', conectando objetivos de negócio e tecnologia através de inovação e boas práticas.',
  authorBio2: 'Profunda expertise em',
  authorBio2b: ', sistemas orientados a eventos, padrões de Microsserviços e arquiteturas de alta escala. Experiência com Data Mesh, Kafka, EKS e AWS Well-Architected Framework.',
  authorBio3: 'Abordagem de desenvolvimento security-first com ampla experiência em CI/CD, Infraestrutura como Código (Terraform, CDK), GitOps e Arquitetura Zero Trust. Conformidade com PCI-DSS, ISO 27001 e GDPR.',
  footerFork: 'diagram-as-code · Fork por Fernando Azevedo',

  // Builder page
  builderLink: 'Construtor',
  builderTitle: 'Construtor YAML',
  builderDesc: 'Construa seu diagrama AWS visualmente — sem precisar escrever YAML',
  defFileSection: 'Arquivo de Definição',
  officialUrl: 'URL Oficial (recomendado)',
  customUrlLabel: 'URL Personalizada',
  localFileLabel: 'Arquivo Local',
  filePathLabel: 'Caminho do arquivo',
  canvasSection: 'Configurações do Canvas',
  canvasDir: 'Direção do Layout',
  resourcesSection: 'Recursos',
  addResource: '+ Adicionar Recurso',
  resourceName: 'Nome (chave YAML)',
  resourceType: 'Tipo',
  resourceTitle: 'Título (opcional)',
  resourcePreset: 'Preset (opcional)',
  resourceDirection: 'Direção (opcional)',
  resourceAlign: 'Alinhamento (opcional)',
  resourceChildren: 'Filhos',
  borderChildrenLabel: 'Filhos na Borda',
  addBorderChild: '+ Adicionar filho na borda',
  positionLabel: 'Posição',
  removeLabel: 'Remover',
  linksSection: 'Conexões',
  addLink: '+ Adicionar Conexão',
  linkSource: 'Origem',
  linkTarget: 'Destino',
  sourcePosLabel: 'Posição de Origem (opcional)',
  targetPosLabel: 'Posição de Destino (opcional)',
  arrowHeadLabel: 'Ponta de Seta',
  edgeLabels: 'Rótulos da Aresta (opcional)',
  yamlPreview: 'Prévia do YAML',
  copyYaml: 'Copiar YAML',
  useInEditor: 'Usar no Editor',
  copied: 'Copiado!',
  noResources: 'Nenhum recurso ainda. Clique em "+ Adicionar Recurso" para começar.',
  noLinks: 'Nenhuma conexão ainda.',
  searchTypes: 'Buscar tipos…',
  selectOptional: '(nenhum)',
  openArrow: 'Aberta',
  noArrow: 'Nenhuma',
  defaultArrow: 'Padrão',
  iconFillLabel: 'Preenchimento do Ícone (opcional)',
  resourceNamePlaceholder: 'ex: MinhaVPC',
  customUrlPlaceholder: 'https://...',
  localFilePlaceholder: '~/minha-definicao.yaml',
  linesCount: (n: number) => `${n} linhas`,
}

export const translations = { en, pt }

// ── Context ───────────────────────────────────────────────────────────────────

type LanguageContextType = {
  lang: Lang
  t: typeof en
  toggle: () => void
}

const LanguageContext = createContext<LanguageContextType>({
  lang: 'en',
  t: en,
  toggle: () => {},
})

function resolveLanguage(): Lang {
  if (typeof window === 'undefined') return 'en'
  const stored = localStorage.getItem('lang')
  if (stored === 'en' || stored === 'pt') return stored
  return navigator.language.toLowerCase().startsWith('pt') ? 'pt' : 'en'
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [lang, setLang] = useState<Lang>(resolveLanguage)

  function toggle() {
    setLang((prev) => {
      const next = prev === 'en' ? 'pt' : 'en'
      localStorage.setItem('lang', next)
      return next
    })
  }

  return (
    <LanguageContext.Provider value={{ lang, t: translations[lang], toggle }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  return useContext(LanguageContext)
}
