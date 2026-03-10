'use client'

import { useState, useCallback, useEffect, useRef } from 'react'
import dynamic from 'next/dynamic'
import Link from 'next/link'
import { ChevronDown, Zap, Github, BookOpen, Wrench } from 'lucide-react'
import DiagramPreview from '@/components/DiagramPreview'
import LanguageSwitcher from '@/components/LanguageSwitcher'
import ThemeSwitcher from '@/components/ThemeSwitcher'
import Tour, { TourStep } from '@/components/Tour'
import { useLanguage } from '@/lib/i18n'

const YamlEditor = dynamic(() => import('@/components/YamlEditor'), {
  ssr: false,
  loading: () => (
    <div className="flex items-center justify-center h-full text-[var(--text-6)] text-sm">
      Loading editor…
    </div>
  ),
})

// ── Example templates ────────────────────────────────────────────────────────

const DEF_URL = 'https://raw.githubusercontent.com/fernandofatech/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml'

const EXAMPLES: Record<string, string> = {
  'BFF Architecture': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"

  Resources:

    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - Clients
        - AWSCloud

    Clients:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - WebClient
        - MobileClient

    WebClient:
      Type: AWS::Diagram::Resource
      Preset: Client
      Title: Web Browser

    MobileClient:
      Type: AWS::Diagram::Resource
      Preset: Mobile client
      Title: Mobile App

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - EdgeLayer
        - AuthLayer
        - BFFLayer
        - VPC

    EdgeLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Route53
        - WAF
        - CDN

    Route53:
      Type: AWS::Route53::HostedZone
      Title: Route 53

    WAF:
      Type: AWS::WAFv2::WebACL
      Title: AWS WAF

    CDN:
      Type: AWS::CloudFront
      Title: CloudFront

    AuthLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Cognito

    Cognito:
      Type: AWS::Cognito::UserPool
      Title: Cognito User Pool

    BFFLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - WebBFFGroup
        - MobileBFFGroup

    WebBFFGroup:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Title: Web BFF
      Direction: vertical
      Children:
        - WebAPIGateway
        - WebBFFLambda

    WebAPIGateway:
      Type: AWS::ApiGateway
      Title: API Gateway (Web)

    WebBFFLambda:
      Type: AWS::Lambda::Function
      Title: Web BFF

    MobileBFFGroup:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Title: Mobile BFF
      Direction: vertical
      Children:
        - MobileAPIGateway
        - MobileBFFLambda

    MobileAPIGateway:
      Type: AWS::ApiGateway
      Title: API Gateway (Mobile)

    MobileBFFLambda:
      Type: AWS::Lambda::Function
      Title: Mobile BFF

    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - InternalALB
        - ServicesStack
        - DataStack
        - MessagingStack

    InternalALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
      Title: Internal ALB

    ServicesStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ProductService
        - OrderService
        - UserService

    ProductService:
      Type: AWS::ECS::Service
      Title: Product Service

    OrderService:
      Type: AWS::ECS::Service
      Title: Order Service

    UserService:
      Type: AWS::ECS::Service
      Title: User Service

    DataStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - RDS
        - Cache
        - NoSQL

    RDS:
      Type: AWS::RDS::DBCluster
      Title: Aurora PostgreSQL

    Cache:
      Type: AWS::ElastiCache::CacheCluster
      Title: ElastiCache Redis

    NoSQL:
      Type: AWS::DynamoDB::Table
      Title: DynamoDB

    MessagingStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Queue
        - Topic

    Queue:
      Type: AWS::SQS::Queue
      Title: SQS Queue

    Topic:
      Type: AWS::SNS::Topic
      Title: SNS Topic

  Links:
    - Source: WebClient
      SourcePosition: E
      Target: Route53
      TargetPosition: W
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTPS"
    - Source: MobileClient
      SourcePosition: E
      Target: Route53
      TargetPosition: W
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTPS"
    - Source: Route53
      Target: WAF
      TargetArrowHead:
        Type: Open
    - Source: WAF
      Target: CDN
      TargetArrowHead:
        Type: Open
    - Source: CDN
      SourcePosition: SSW
      Target: WebAPIGateway
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        TargetLeft:
          Title: "Web"
    - Source: CDN
      SourcePosition: SSE
      Target: MobileAPIGateway
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        TargetRight:
          Title: "Mobile"
    - Source: WebBFFLambda
      Target: Cognito
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "JWT"
    - Source: MobileBFFLambda
      Target: Cognito
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "JWT"
    - Source: WebBFFLambda
      SourcePosition: S
      Target: InternalALB
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: MobileBFFLambda
      SourcePosition: S
      Target: InternalALB
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: InternalALB
      Target: ProductService
      TargetArrowHead:
        Type: Open
    - Source: InternalALB
      Target: OrderService
      TargetArrowHead:
        Type: Open
    - Source: InternalALB
      Target: UserService
      TargetArrowHead:
        Type: Open
    - Source: ProductService
      Target: RDS
      TargetArrowHead:
        Type: Open
    - Source: OrderService
      Target: RDS
      TargetArrowHead:
        Type: Open
    - Source: UserService
      Target: Cache
      TargetArrowHead:
        Type: Open
    - Source: OrderService
      Target: NoSQL
      TargetArrowHead:
        Type: Open
    - Source: OrderService
      Target: Queue
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "async"
    - Source: Queue
      Target: Topic
      TargetArrowHead:
        Type: Open
`,

  'Event-Driven (SNS/SQS/Lambda)': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - ClientApp
        - AWSCloud

    ClientApp:
      Type: AWS::Diagram::Resource
      Preset: Client
      Title: Client App

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - APILayer
        - MessagingLayer
        - ProcessingLayer
        - StorageLayer

    APILayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - APIGW
        - EventLambda

    APIGW:
      Type: AWS::ApiGateway
      Title: API Gateway

    EventLambda:
      Type: AWS::Lambda::Function
      Title: Event Publisher

    MessagingLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - OrdersTopic
        - OrdersQueue
        - PaymentsQueue
        - NotifQueue

    OrdersTopic:
      Type: AWS::SNS::Topic
      Title: Orders Topic

    OrdersQueue:
      Type: AWS::SQS::Queue
      Title: Orders Queue

    PaymentsQueue:
      Type: AWS::SQS::Queue
      Title: Payments Queue

    NotifQueue:
      Type: AWS::SQS::Queue
      Title: Notifications Queue

    ProcessingLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - OrderProcessor
        - PaymentProcessor
        - NotifProcessor

    OrderProcessor:
      Type: AWS::Lambda::Function
      Title: Order Processor

    PaymentProcessor:
      Type: AWS::Lambda::Function
      Title: Payment Processor

    NotifProcessor:
      Type: AWS::Lambda::Function
      Title: Notif Sender

    StorageLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - OrdersTable
        - PaymentsTable
        - AuditBucket

    OrdersTable:
      Type: AWS::DynamoDB::Table
      Title: Orders

    PaymentsTable:
      Type: AWS::DynamoDB::Table
      Title: Payments

    AuditBucket:
      Type: AWS::S3::Bucket
      Title: Audit Logs

  Links:
    - Source: ClientApp
      Target: APIGW
      TargetArrowHead:
        Type: Open
    - Source: APIGW
      Target: EventLambda
      TargetArrowHead:
        Type: Open
    - Source: EventLambda
      Target: OrdersTopic
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "publish"
    - Source: OrdersTopic
      SourcePosition: SSW
      Target: OrdersQueue
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: OrdersTopic
      SourcePosition: S
      Target: PaymentsQueue
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: OrdersTopic
      SourcePosition: SSE
      Target: NotifQueue
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: OrdersQueue
      Target: OrderProcessor
      TargetArrowHead:
        Type: Open
    - Source: PaymentsQueue
      Target: PaymentProcessor
      TargetArrowHead:
        Type: Open
    - Source: NotifQueue
      Target: NotifProcessor
      TargetArrowHead:
        Type: Open
    - Source: OrderProcessor
      Target: OrdersTable
      TargetArrowHead:
        Type: Open
    - Source: PaymentProcessor
      Target: PaymentsTable
      TargetArrowHead:
        Type: Open
    - Source: NotifProcessor
      Target: AuditBucket
      TargetArrowHead:
        Type: Open
`,

  'Serverless REST API': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - MobileUser
        - AWSCloud

    MobileUser:
      Type: AWS::Diagram::Resource
      Preset: Mobile client
      Title: Mobile / Web App

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - EdgeLayer
        - AuthLayer
        - ComputeLayer
        - DataLayer

    EdgeLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Route53
        - CDN
        - APIGW

    Route53:
      Type: AWS::Route53::HostedZone
      Title: Route 53

    CDN:
      Type: AWS::CloudFront
      Title: CloudFront

    APIGW:
      Type: AWS::ApiGatewayV2::Api
      Title: HTTP API

    AuthLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - UserPool

    UserPool:
      Type: AWS::Cognito::UserPool
      Title: Cognito User Pool

    ComputeLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - UsersLambda
        - ProductsLambda
        - OrdersLambda

    UsersLambda:
      Type: AWS::Lambda::Function
      Title: Users fn

    ProductsLambda:
      Type: AWS::Lambda::Function
      Title: Products fn

    OrdersLambda:
      Type: AWS::Lambda::Function
      Title: Orders fn

    DataLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - UsersTable
        - ProductsTable
        - OrdersTable
        - Cache

    UsersTable:
      Type: AWS::DynamoDB::Table
      Title: Users

    ProductsTable:
      Type: AWS::DynamoDB::Table
      Title: Products

    OrdersTable:
      Type: AWS::DynamoDB::Table
      Title: Orders

    Cache:
      Type: AWS::ElastiCache::CacheCluster
      Title: Redis Cache

  Links:
    - Source: MobileUser
      Target: Route53
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTPS"
    - Source: Route53
      Target: CDN
      TargetArrowHead:
        Type: Open
    - Source: CDN
      Target: APIGW
      TargetArrowHead:
        Type: Open
    - Source: APIGW
      Target: UserPool
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "auth"
    - Source: APIGW
      SourcePosition: SSW
      Target: UsersLambda
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: APIGW
      SourcePosition: S
      Target: ProductsLambda
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: APIGW
      SourcePosition: SSE
      Target: OrdersLambda
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: UsersLambda
      Target: UsersTable
      TargetArrowHead:
        Type: Open
    - Source: ProductsLambda
      Target: ProductsTable
      TargetArrowHead:
        Type: Open
    - Source: OrdersLambda
      Target: OrdersTable
      TargetArrowHead:
        Type: Open
    - Source: UsersLambda
      Target: Cache
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "cache"
`,

  'MFE with CloudFront': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - Browser
        - AWSCloud

    Browser:
      Type: AWS::Diagram::Resource
      Preset: Client
      Title: Browser

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - CDNLayer
        - StaticLayer
        - AuthLayer
        - APILayer
        - DataLayer

    CDNLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - CloudFront

    CloudFront:
      Type: AWS::CloudFront
      Title: CloudFront CDN

    StaticLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ShellBucket
        - HomeBucket
        - CartBucket
        - AssetsBucket

    ShellBucket:
      Type: AWS::S3::Bucket
      Title: Shell App

    HomeBucket:
      Type: AWS::S3::Bucket
      Title: Home MFE

    CartBucket:
      Type: AWS::S3::Bucket
      Title: Cart MFE

    AssetsBucket:
      Type: AWS::S3::Bucket
      Title: Static Assets

    AuthLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Cognito

    Cognito:
      Type: AWS::Cognito::UserPool
      Title: Cognito

    APILayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - APIGW
        - HomeAPILambda
        - CartAPILambda

    APIGW:
      Type: AWS::ApiGatewayV2::Api
      Title: API Gateway

    HomeAPILambda:
      Type: AWS::Lambda::Function
      Title: Home API

    CartAPILambda:
      Type: AWS::Lambda::Function
      Title: Cart API

    DataLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ProductsDB
        - CartDB

    ProductsDB:
      Type: AWS::DynamoDB::Table
      Title: Products

    CartDB:
      Type: AWS::DynamoDB::Table
      Title: Cart

  Links:
    - Source: Browser
      Target: CloudFront
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTPS"
    - Source: CloudFront
      SourcePosition: SSW
      Target: ShellBucket
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        TargetLeft:
          Title: "/"
    - Source: CloudFront
      SourcePosition: S
      Target: HomeBucket
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        TargetLeft:
          Title: "/home"
    - Source: CloudFront
      SourcePosition: SSE
      Target: CartBucket
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        TargetRight:
          Title: "/cart"
    - Source: CloudFront
      SourcePosition: E
      Target: APIGW
      TargetPosition: W
      TargetArrowHead:
        Type: Open
      Labels:
        TargetRight:
          Title: "/api"
    - Source: APIGW
      Target: Cognito
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "JWT"
    - Source: APIGW
      SourcePosition: SSW
      Target: HomeAPILambda
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: APIGW
      SourcePosition: SSE
      Target: CartAPILambda
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: HomeAPILambda
      Target: ProductsDB
      TargetArrowHead:
        Type: Open
    - Source: CartAPILambda
      Target: CartDB
      TargetArrowHead:
        Type: Open
`,

  'ECS Microservices': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - User
        - AWSCloud

    User:
      Type: AWS::Diagram::Resource
      Preset: User

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC

    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - PublicSubnet
        - PrivateSubnet
        - DataSubnet
      BorderChildren:
        - Position: S
          Resource: IGW

    IGW:
      Type: AWS::EC2::InternetGateway

    PublicSubnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Title: Public Subnet
      Children:
        - ALB

    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer

    PrivateSubnet:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Title: App Subnet (ECS)
      Children:
        - ServicesStack

    ServicesStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - UsersService
        - ProductsService
        - OrdersService

    UsersService:
      Type: AWS::ECS::Service
      Title: Users Service

    ProductsService:
      Type: AWS::ECS::Service
      Title: Products Service

    OrdersService:
      Type: AWS::ECS::Service
      Title: Orders Service

    DataSubnet:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Title: Data Subnet
      Children:
        - DataStack

    DataStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - AuroraDB
        - RedisCache
        - EventsQueue

    AuroraDB:
      Type: AWS::RDS::DBCluster
      Title: Aurora PostgreSQL

    RedisCache:
      Type: AWS::ElastiCache::CacheCluster
      Title: ElastiCache Redis

    EventsQueue:
      Type: AWS::SQS::Queue
      Title: Events Queue

  Links:
    - Source: User
      SourcePosition: N
      Target: IGW
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: IGW
      SourcePosition: N
      Target: ALB
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: NNW
      Target: UsersService
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: N
      Target: ProductsService
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: NNE
      Target: OrdersService
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: UsersService
      Target: AuroraDB
      TargetArrowHead:
        Type: Open
    - Source: ProductsService
      Target: AuroraDB
      TargetArrowHead:
        Type: Open
    - Source: OrdersService
      Target: AuroraDB
      TargetArrowHead:
        Type: Open
    - Source: UsersService
      Target: RedisCache
      TargetArrowHead:
        Type: Open
    - Source: OrdersService
      Target: EventsQueue
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "async"
`,

  'CI/CD Pipeline': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - DevLayer
        - AWSCloud

    DevLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Developer

    Developer:
      Type: AWS::Diagram::Resource
      Preset: User
      Title: Developer

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - PipelineLayer
        - BuildLayer
        - RegistryLayer
        - DeployLayer

    PipelineLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Pipeline

    Pipeline:
      Type: AWS::CodePipeline::Pipeline
      Title: CodePipeline

    BuildLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Build

    Build:
      Type: AWS::CodeBuild::Project
      Title: CodeBuild

    RegistryLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Registry

    Registry:
      Type: AWS::ECR
      Title: ECR

    DeployLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ECSCluster
        - APIService
        - WorkerService

    ECSCluster:
      Type: AWS::ECS::Cluster
      Title: ECS Cluster

    APIService:
      Type: AWS::ECS::Service
      Title: API Service

    WorkerService:
      Type: AWS::ECS::Service
      Title: Worker Service

  Links:
    - Source: Developer
      Target: Pipeline
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "git push"
    - Source: Pipeline
      Target: Build
      TargetArrowHead:
        Type: Open
    - Source: Build
      Target: Registry
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "docker push"
    - Source: Registry
      SourcePosition: SSW
      Target: APIService
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: Registry
      SourcePosition: SSE
      Target: WorkerService
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: Pipeline
      SourcePosition: SSE
      Target: APIService
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "deploy"
`,

  'Data Lake & Analytics': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - Sources
        - AWSCloud

    Sources:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - AppSource
        - StreamSource

    AppSource:
      Type: AWS::Diagram::Resource
      Preset: Client
      Title: Applications

    StreamSource:
      Type: AWS::Diagram::Resource
      Preset: Mobile client
      Title: IoT / Events

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - IngestionLayer
        - StorageLayer
        - ProcessingLayer
        - AnalyticsLayer

    IngestionLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - KinesisStream
        - APIGW

    KinesisStream:
      Type: AWS::Kinesis::Stream
      Title: Kinesis Data Stream

    APIGW:
      Type: AWS::ApiGateway
      Title: API Gateway

    StorageLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - RawBucket
        - ProcessedBucket
        - CuratedBucket

    RawBucket:
      Type: AWS::S3::Bucket
      Title: Raw Zone

    ProcessedBucket:
      Type: AWS::S3::Bucket
      Title: Processed Zone

    CuratedBucket:
      Type: AWS::S3::Bucket
      Title: Curated Zone

    ProcessingLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - GlueETL
        - TransformLambda

    GlueETL:
      Type: AWS::Glue
      Title: AWS Glue ETL

    TransformLambda:
      Type: AWS::Lambda::Function
      Title: Transform Lambda

    AnalyticsLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - AthenaQuery
        - EventsScheduler

    AthenaQuery:
      Type: AWS::Athena
      Title: Amazon Athena

    EventsScheduler:
      Type: AWS::Events::Rule
      Title: EventBridge Scheduler

  Links:
    - Source: StreamSource
      Target: KinesisStream
      TargetArrowHead:
        Type: Open
    - Source: AppSource
      Target: APIGW
      TargetArrowHead:
        Type: Open
    - Source: KinesisStream
      Target: RawBucket
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "stream"
    - Source: APIGW
      Target: RawBucket
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "batch"
    - Source: RawBucket
      Target: GlueETL
      TargetArrowHead:
        Type: Open
    - Source: GlueETL
      Target: ProcessedBucket
      TargetArrowHead:
        Type: Open
    - Source: ProcessedBucket
      Target: TransformLambda
      TargetArrowHead:
        Type: Open
    - Source: TransformLambda
      Target: CuratedBucket
      TargetArrowHead:
        Type: Open
    - Source: CuratedBucket
      Target: AthenaQuery
      TargetArrowHead:
        Type: Open
    - Source: EventsScheduler
      Target: GlueETL
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "schedule"
`,

  'EKS Platform': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - DevTeam
        - AWSCloud

    DevTeam:
      Type: AWS::Diagram::Resource
      Preset: User
      Title: Dev Team

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - EdgeLayer
        - RegistryLayer
        - VPC

    EdgeLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ALBIngress

    ALBIngress:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
      Title: ALB Ingress

    RegistryLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ECR

    ECR:
      Type: AWS::ECR
      Title: ECR

    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - PublicSubnet
        - PrivateSubnet
        - DataSubnet
      BorderChildren:
        - Position: S
          Resource: IGW

    IGW:
      Type: AWS::EC2::InternetGateway

    PublicSubnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Title: Public Subnet
      Children:
        - NAT

    NAT:
      Type: AWS::EC2::NatGateway

    PrivateSubnet:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Title: EKS Node Group
      Children:
        - EKSCluster

    EKSCluster:
      Type: AWS::EKS
      Title: EKS Cluster

    DataSubnet:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Title: Data Subnet
      Children:
        - DataStack

    DataStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Database
        - CacheCluster

    Database:
      Type: AWS::RDS::DBCluster
      Title: Aurora PostgreSQL

    CacheCluster:
      Type: AWS::ElastiCache::CacheCluster
      Title: Redis

  Links:
    - Source: DevTeam
      SourcePosition: E
      Target: ECR
      TargetPosition: W
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "docker push"
    - Source: DevTeam
      SourcePosition: S
      Target: ALBIngress
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "kubectl"
    - Source: ALBIngress
      Target: EKSCluster
      TargetArrowHead:
        Type: Open
    - Source: ECR
      Target: EKSCluster
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "pull image"
    - Source: EKSCluster
      Target: Database
      TargetArrowHead:
        Type: Open
    - Source: EKSCluster
      Target: CacheCluster
      TargetArrowHead:
        Type: Open
    - Source: NAT
      Target: IGW
      TargetArrowHead:
        Type: Open
`,

  'Multi-tier Web App': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - User
        - AWSCloud

    User:
      Type: AWS::Diagram::Resource
      Preset: User

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - EdgeLayer
        - VPC

    EdgeLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - DNS
        - CDN
        - WAF

    DNS:
      Type: AWS::Route53::HostedZone
      Title: Route 53

    CDN:
      Type: AWS::CloudFront
      Title: CloudFront

    WAF:
      Type: AWS::WAFv2::WebACL
      Title: AWS WAF

    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - PublicSubnet
        - AppSubnet
        - DBSubnet
      BorderChildren:
        - Position: S
          Resource: IGW

    IGW:
      Type: AWS::EC2::InternetGateway

    PublicSubnet:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Title: Public Subnet
      Children:
        - ALB

    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer

    AppSubnet:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Title: App Subnet
      Children:
        - ASG

    ASG:
      Type: AWS::AutoScaling::AutoScalingGroup
      Children:
        - AppServer1
        - AppServer2

    AppServer1:
      Type: AWS::EC2::Instance
      Title: App Server

    AppServer2:
      Type: AWS::EC2::Instance
      Title: App Server

    DBSubnet:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Title: Data Subnet
      Children:
        - DBStack

    DBStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - PrimaryDB
        - ReplicaDB
        - SessionCache

    PrimaryDB:
      Type: AWS::RDS::DBCluster
      Title: Aurora (Primary)

    ReplicaDB:
      Type: AWS::RDS::DBInstance
      Title: Aurora (Replica)

    SessionCache:
      Type: AWS::ElastiCache::CacheCluster
      Title: ElastiCache Redis

  Links:
    - Source: User
      Target: DNS
      TargetArrowHead:
        Type: Open
    - Source: DNS
      Target: CDN
      TargetArrowHead:
        Type: Open
    - Source: CDN
      Target: WAF
      TargetArrowHead:
        Type: Open
    - Source: WAF
      Target: ALB
      TargetArrowHead:
        Type: Open
    - Source: IGW
      SourcePosition: N
      Target: ALB
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: NNW
      Target: AppServer1
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: NNE
      Target: AppServer2
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: AppServer1
      Target: PrimaryDB
      TargetArrowHead:
        Type: Open
    - Source: AppServer2
      Target: PrimaryDB
      TargetArrowHead:
        Type: Open
    - Source: AppServer1
      Target: SessionCache
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "session"
    - Source: PrimaryDB
      Target: ReplicaDB
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "replication"
`,

  'Step Functions Saga': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "${DEF_URL}"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - ClientLayer
        - AWSCloud

    ClientLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Client

    Client:
      Type: AWS::Diagram::Resource
      Preset: Mobile client
      Title: Client App

    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - APILayer
        - OrchestratorLayer
        - StepsLayer
        - CompensationLayer
        - StorageLayer

    APILayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - APIGW

    APIGW:
      Type: AWS::ApiGateway
      Title: API Gateway

    OrchestratorLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - StateMachine

    StateMachine:
      Type: AWS::StepFunctions
      Title: Order Saga

    StepsLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ReserveInventoryFn
        - ProcessPaymentFn
        - ShipOrderFn

    ReserveInventoryFn:
      Type: AWS::Lambda::Function
      Title: Reserve Inventory

    ProcessPaymentFn:
      Type: AWS::Lambda::Function
      Title: Process Payment

    ShipOrderFn:
      Type: AWS::Lambda::Function
      Title: Ship Order

    CompensationLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - ReleaseInventoryFn
        - RefundPaymentFn
        - NotifyFailureFn

    ReleaseInventoryFn:
      Type: AWS::Lambda::Function
      Title: Release Inventory

    RefundPaymentFn:
      Type: AWS::Lambda::Function
      Title: Refund Payment

    NotifyFailureFn:
      Type: AWS::Lambda::Function
      Title: Notify Failure

    StorageLayer:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - OrdersTable
        - InventoryTable
        - DLQ

    OrdersTable:
      Type: AWS::DynamoDB::Table
      Title: Orders

    InventoryTable:
      Type: AWS::DynamoDB::Table
      Title: Inventory

    DLQ:
      Type: AWS::SQS::Queue
      Title: Dead Letter Queue

  Links:
    - Source: Client
      Target: APIGW
      TargetArrowHead:
        Type: Open
    - Source: APIGW
      Target: StateMachine
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "StartExecution"
    - Source: StateMachine
      SourcePosition: SSW
      Target: ReserveInventoryFn
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: StateMachine
      SourcePosition: S
      Target: ProcessPaymentFn
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: StateMachine
      SourcePosition: SSE
      Target: ShipOrderFn
      TargetPosition: N
      TargetArrowHead:
        Type: Open
    - Source: StateMachine
      SourcePosition: SW
      Target: ReleaseInventoryFn
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "compensate"
    - Source: StateMachine
      SourcePosition: S
      Target: RefundPaymentFn
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "compensate"
    - Source: StateMachine
      SourcePosition: SE
      Target: NotifyFailureFn
      TargetPosition: N
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "compensate"
    - Source: ReserveInventoryFn
      Target: InventoryTable
      TargetArrowHead:
        Type: Open
    - Source: ProcessPaymentFn
      Target: OrdersTable
      TargetArrowHead:
        Type: Open
    - Source: ShipOrderFn
      Target: OrdersTable
      TargetArrowHead:
        Type: Open
    - Source: NotifyFailureFn
      Target: DLQ
      TargetArrowHead:
        Type: Open
`,

  'ALB + EC2': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/fernandofatech/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
        - User
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - VPCPublicStack
        - ALB
      BorderChildren:
        - Position: S
          Resource: IGW
    VPCPublicStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - VPCPublicSubnet1
        - VPCPublicSubnet2
    VPCPublicSubnet1:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet1Instance
    VPCPublicSubnet1Instance:
      Type: AWS::EC2::Instance
    VPCPublicSubnet2:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet2Instance
    VPCPublicSubnet2Instance:
      Type: AWS::EC2::Instance
    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGW:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    User:
      Type: AWS::Diagram::Resource
      Preset: User
  Links:
    - Source: ALB
      SourcePosition: NNW
      Target: VPCPublicSubnet1Instance
      TargetPosition: SSE
      TargetArrowHead:
        Type: Open
    - Source: ALB
      SourcePosition: NNE
      Target: VPCPublicSubnet2Instance
      TargetPosition: SSW
      TargetArrowHead:
        Type: Open
    - Source: IGW
      SourcePosition: N
      Target: ALB
      TargetPosition: S
      TargetArrowHead:
        Type: Open
    - Source: User
      SourcePosition: N
      Target: IGW
      TargetPosition: S
      TargetArrowHead:
        Type: Open
`,

  'VPC + NAT Gateway': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/fernandofatech/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
        - User
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - VPCPublicSubnetStack
        - VPCPrivateSubnetStack
    VPCPublicSubnetStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - VPCPublicSubnet1
        - VPCPublicSubnet2
    VPCPublicSubnet1:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet1NatGateway
    VPCPublicSubnet1NatGateway:
      Type: AWS::EC2::NatGateway
    VPCPublicSubnet2:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - VPCPublicSubnet2NatGateway
    VPCPublicSubnet2NatGateway:
      Type: AWS::EC2::NatGateway
    VPCPrivateSubnetStack:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - VPCPrivateSubnet1
        - VPCPrivateSubnet2
    VPCPrivateSubnet1:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Children:
        - VPCPrivateSubnet1Instance
    VPCPrivateSubnet1Instance:
      Type: AWS::EC2::Instance
    VPCPrivateSubnet2:
      Type: AWS::EC2::Subnet
      Preset: PrivateSubnet
      Children:
        - VPCPrivateSubnet2Instance
    VPCPrivateSubnet2Instance:
      Type: AWS::EC2::Instance
    User:
      Type: AWS::Diagram::Resource
      Preset: User
`,

  'ALB + Auto Scaling': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/fernandofatech/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
        - User
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - VPC
    VPC:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - AutoScalingGroup
        - ALB
      BorderChildren:
        - Position: S
          Resource: IGW
    AutoScalingGroup:
      Type: AWS::AutoScaling::AutoScalingGroup
      Children:
        - Instance1
        - Instance2
    Instance1:
      Type: AWS::EC2::Instance
    Instance2:
      Type: AWS::EC2::Instance
    ALB:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGW:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    User:
      Type: AWS::Diagram::Resource
      Preset: User
  Links:
    - Source: ALB
      SourcePosition: NNW
      Target: Instance1
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
    - Source: ALB
      SourcePosition: NNE
      Target: Instance2
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceRight:
          Title: "HTTP:80"
    - Source: IGW
      SourcePosition: N
      Target: ALB
      TargetPosition: S
      TargetArrowHead:
        Type: Open
      Labels:
        SourceLeft:
          Title: "HTTP:80"
        SourceRight:
          Title: "HTTPS:443"
    - Source: User
      SourcePosition: N
      Target: IGW
      TargetPosition: S
      TargetArrowHead:
        Type: Open
`,

  'Multi-Region': `Diagram:
  DefinitionFiles:
    - Type: URL
      Url: "https://raw.githubusercontent.com/fernandofatech/diagram-as-code/main/definitions/definition-for-aws-icons-light.yaml"
  Resources:
    Canvas:
      Type: AWS::Diagram::Canvas
      Direction: vertical
      Children:
        - AWSCloud
        - User
    AWSCloud:
      Type: AWS::Diagram::Cloud
      Direction: vertical
      Preset: AWSCloudNoLogo
      Align: center
      Children:
        - Regions
    Regions:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - UsEast1
        - UsWest2
    UsEast1:
      Type: AWS::Region
      Title: us-east-1
      Direction: vertical
      Children:
        - VpcEast
    VpcEast:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - SubnetStackEast
        - ALBEast
      BorderChildren:
        - Position: S
          Resource: IGWEast
    SubnetStackEast:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Subnet1East
        - Subnet2East
    Subnet1East:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance1East
    Instance1East:
      Type: AWS::EC2::Instance
    Subnet2East:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance2East
    Instance2East:
      Type: AWS::EC2::Instance
    ALBEast:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGWEast:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    UsWest2:
      Type: AWS::Region
      Title: us-west-2
      Direction: vertical
      Children:
        - VpcWest
    VpcWest:
      Type: AWS::EC2::VPC
      Direction: vertical
      Children:
        - SubnetStackWest
        - ALBWest
      BorderChildren:
        - Position: S
          Resource: IGWWest
    SubnetStackWest:
      Type: AWS::Diagram::HorizontalStack
      Children:
        - Subnet1West
        - Subnet2West
    Subnet1West:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance1West
    Instance1West:
      Type: AWS::EC2::Instance
    Subnet2West:
      Type: AWS::EC2::Subnet
      Preset: PublicSubnet
      Children:
        - Instance2West
    Instance2West:
      Type: AWS::EC2::Instance
    ALBWest:
      Type: AWS::ElasticLoadBalancingV2::LoadBalancer
      Preset: Application Load Balancer
    IGWWest:
      Type: AWS::EC2::InternetGateway
      IconFill:
        Type: rect
    User:
      Type: AWS::Diagram::Resource
      Preset: User
  Links:
    - Source: ALBEast
      Target: Instance1East
      TargetArrowHead:
        Type: Open
    - Source: ALBEast
      Target: Instance2East
      TargetArrowHead:
        Type: Open
    - Source: IGWEast
      Target: ALBEast
      TargetArrowHead:
        Type: Open
    - Source: User
      Target: IGWEast
      TargetArrowHead:
        Type: Open
    - Source: ALBWest
      Target: Instance1West
      TargetArrowHead:
        Type: Open
    - Source: ALBWest
      Target: Instance2West
      TargetArrowHead:
        Type: Open
    - Source: IGWWest
      Target: ALBWest
      TargetArrowHead:
        Type: Open
    - Source: User
      Target: IGWWest
      TargetArrowHead:
        Type: Open
`,
}

const EXAMPLE_NAMES = Object.keys(EXAMPLES)

// ── Page component ───────────────────────────────────────────────────────────

const EDITOR_TOUR: TourStep[] = [
  { target: 'editor-logo', title: 'diagram-as-code', description: 'Generate AWS architecture diagrams from YAML. Write code, get pixel-perfect diagrams.', position: 'bottom' },
  { target: 'editor-examples', title: 'Examples', description: '10 real-world templates: BFF, Event-Driven, Serverless API, MFE, ECS, CI/CD, Data Lake, EKS, Multi-tier Web, Step Functions Saga.', position: 'bottom' },
  { target: 'editor-format', title: 'Output Format', description: 'Toggle between PNG, PDF, and draw.io depending on whether you need preview images, printable documents, or editable diagrams.', position: 'bottom' },
  { target: 'editor-generate', title: 'Generate', description: 'Click to render your diagram. You can also press Ctrl+Enter anytime.', position: 'bottom' },
  { target: 'editor-panel', title: 'YAML Editor', description: 'Write your architecture in YAML. The editor has syntax highlighting powered by Monaco (VS Code engine).', position: 'right' },
  { target: 'editor-preview', title: 'Preview', description: 'Your generated diagram appears here. Download PNG, PDF, or draw.io from the preview panel.', position: 'left' },
  { target: 'editor-builder', title: 'Visual Builder', description: 'No YAML experience? Use the Builder to create diagrams with a visual form — no code required.', position: 'bottom' },
]

export default function Home() {
  const { t } = useLanguage()
  const [yaml, setYaml] = useState(EXAMPLES['BFF Architecture'])
  const [format, setFormat] = useState<'png' | 'drawio' | 'pdf'>('png')
  const [imageUrl, setImageUrl] = useState<string | null>(null)
  const [drawioContent, setDrawioContent] = useState<string | null>(null)
  const [pdfUrl, setPdfUrl] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [examplesOpen, setExamplesOpen] = useState(false)
  const prevImageUrl = useRef<string | null>(null)
  const prevPdfUrl = useRef<string | null>(null)

  // Pick up YAML from builder
  useEffect(() => {
    const saved = localStorage.getItem('builder_yaml')
    if (saved) {
      setYaml(saved)
      localStorage.removeItem('builder_yaml')
    }
  }, [])

  useEffect(() => {
    return () => {
      if (prevImageUrl.current) URL.revokeObjectURL(prevImageUrl.current)
      if (prevPdfUrl.current) URL.revokeObjectURL(prevPdfUrl.current)
    }
  }, [])

  // Ctrl+Enter shortcut
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault()
        generate()
      }
    }
    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [yaml, format])

  const generate = useCallback(async () => {
    if (!yaml.trim()) return
    setLoading(true)
    setError(null)

    try {
      const res = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ yaml, format }),
      })

      if (!res.ok) {
        const body = await res.json().catch(() => ({ error: 'Unknown error' }))
        throw new Error(body.error ?? 'Generation failed')
      }

      if (format === 'png') {
        const blob = await res.blob()
        const url = URL.createObjectURL(blob)
        if (prevImageUrl.current) URL.revokeObjectURL(prevImageUrl.current)
        prevImageUrl.current = url
        setImageUrl(url)
        setDrawioContent(null)
        if (prevPdfUrl.current) URL.revokeObjectURL(prevPdfUrl.current)
        prevPdfUrl.current = null
        setPdfUrl(null)
      } else if (format === 'pdf') {
        const blob = await res.blob()
        const url = URL.createObjectURL(blob)
        if (prevPdfUrl.current) URL.revokeObjectURL(prevPdfUrl.current)
        prevPdfUrl.current = url
        setPdfUrl(url)
        setImageUrl(null)
        setDrawioContent(null)
      } else {
        const text = await res.text()
        setDrawioContent(text)
        setImageUrl(null)
        if (prevPdfUrl.current) URL.revokeObjectURL(prevPdfUrl.current)
        prevPdfUrl.current = null
        setPdfUrl(null)
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err))
    } finally {
      setLoading(false)
    }
  }, [yaml, format])

  function loadExample(name: string) {
    setYaml(EXAMPLES[name])
    setExamplesOpen(false)
    if (prevImageUrl.current) URL.revokeObjectURL(prevImageUrl.current)
    prevImageUrl.current = null
    if (prevPdfUrl.current) URL.revokeObjectURL(prevPdfUrl.current)
    prevPdfUrl.current = null
    setImageUrl(null)
    setDrawioContent(null)
    setPdfUrl(null)
    setError(null)
  }

  return (
    <div className="flex flex-col h-screen bg-[var(--bg)] overflow-hidden">

      {/* ── Header ─────────────────────────────────────────────────────────── */}
      <header className="flex items-center justify-between px-4 h-12 border-b border-[var(--border)] flex-shrink-0">
        <div className="flex items-center gap-3">
          {/* Logo */}
          <div className="flex items-center gap-2" data-tour="editor-logo">
            <div className="w-6 h-6 bg-[#FF9900] rounded flex items-center justify-center">
              <svg viewBox="0 0 16 16" fill="white" className="w-4 h-4">
                <rect x="1" y="1" width="6" height="6" rx="1" />
                <rect x="9" y="1" width="6" height="6" rx="1" />
                <rect x="1" y="9" width="6" height="6" rx="1" />
                <rect x="9" y="9" width="6" height="6" rx="1" />
              </svg>
            </div>
            <span className="text-sm font-semibold text-[var(--text)] tracking-tight">
              diagram-as-code
            </span>
          </div>

          {/* Divider */}
          <div className="w-px h-4 bg-[var(--border)]" />

          {/* Examples dropdown */}
          <div className="relative" data-tour="editor-examples">
            <button
              onClick={() => setExamplesOpen((o) => !o)}
              className="flex items-center gap-1.5 text-xs text-[var(--text-3)] hover:text-[var(--text)] transition-colors px-2 py-1 rounded hover:bg-[var(--surface)]"
            >
              {t.examples}
              <ChevronDown size={12} className={`transition-transform ${examplesOpen ? 'rotate-180' : ''}`} />
            </button>
            {examplesOpen && (
              <>
                <div
                  className="fixed inset-0 z-10"
                  onClick={() => setExamplesOpen(false)}
                />
                <div className="absolute top-full left-0 mt-1 bg-[var(--surface)] border border-[var(--border)] rounded-lg shadow-xl z-20 min-w-[180px] py-1 overflow-hidden">
                  {EXAMPLE_NAMES.map((name) => (
                    <button
                      key={name}
                      onClick={() => loadExample(name)}
                      className="w-full text-left px-3 py-2 text-xs text-[var(--text-2)] hover:bg-[var(--surface-hover)] hover:text-[var(--text)] transition-colors"
                    >
                      {name}
                    </button>
                  ))}
                </div>
              </>
            )}
          </div>
        </div>

        {/* Right controls */}
        <div className="flex items-center gap-2">
          {/* Language switcher */}
          <LanguageSwitcher />

          {/* Theme switcher */}
          <ThemeSwitcher />

          {/* Format toggle */}
          <div data-tour="editor-format" className="flex items-center bg-[var(--surface)] border border-[var(--border)] rounded-md p-0.5 text-xs">
            {(['png', 'pdf', 'drawio'] as const).map((f) => (
              <button
                key={f}
                onClick={() => setFormat(f)}
                className={`px-3 py-1 rounded transition-all ${
                  format === f
                    ? 'bg-[#FF9900] text-[var(--accent-contrast)] font-semibold'
                    : 'text-[var(--text-3)] hover:text-[var(--text-2)]'
                }`}
              >
                {f === 'png' ? 'PNG' : f === 'pdf' ? 'PDF' : 'draw.io'}
              </button>
            ))}
          </div>

          {/* Generate button */}
          <button
            data-tour="editor-generate"
            onClick={generate}
            disabled={loading || !yaml.trim()}
            className="flex items-center gap-1.5 px-4 py-1.5 bg-[#FF9900] hover:bg-[#ffb340] disabled:opacity-40 disabled:cursor-not-allowed text-[var(--accent-contrast)] text-xs font-semibold rounded-md transition-colors"
          >
            {loading ? (
              <div className="w-3.5 h-3.5 border-2 border-[var(--accent-contrast)] border-t-transparent rounded-full animate-spin" />
            ) : (
              <Zap size={13} />
            )}
            {t.generate}
          </button>

          {/* Builder link */}
          <Link
            data-tour="editor-builder"
            href="/builder"
            className="flex items-center gap-1.5 text-xs text-[var(--text-5)] hover:text-[var(--text-3)] transition-colors px-2 py-1 rounded hover:bg-[var(--surface)]"
          >
            <Wrench size={13} />
            {t.builderLink}
          </Link>

          {/* Docs link */}
          <Link
            href="/docs"
            className="flex items-center gap-1.5 text-xs text-[var(--text-5)] hover:text-[var(--text-3)] transition-colors px-2 py-1 rounded hover:bg-[var(--surface)]"
          >
            <BookOpen size={13} />
            {t.docs}
          </Link>

          {/* GitHub link */}
          <a
            href="https://github.com/fernandofatech/diagram-as-code"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[var(--text-5)] hover:text-[var(--text-3)] transition-colors p-1"
            aria-label="GitHub"
          >
            <Github size={16} />
          </a>
        </div>
      </header>

      {/* ── Main content ────────────────────────────────────────────────────── */}
      <main className="flex flex-1 overflow-hidden">
        {/* Editor panel */}
        <div data-tour="editor-panel" className="w-1/2 flex flex-col overflow-hidden border-r border-[var(--border)]">
          <div className="flex items-center justify-between px-4 py-2 border-b border-[var(--border)] flex-shrink-0">
            <span className="text-xs text-[var(--text-5)] font-medium uppercase tracking-wider">
              {t.yamlEditor}
            </span>
            <span className="text-xs text-[var(--text-6)]">
              {t.lines(yaml.split('\n').length)}
            </span>
          </div>
          <div className="flex-1 overflow-hidden">
            <YamlEditor value={yaml} onChange={setYaml} />
          </div>
        </div>

        {/* Preview panel */}
        <div data-tour="editor-preview" className="w-1/2 flex flex-col overflow-hidden">
          <div className="flex items-center px-4 py-2 border-b border-[var(--border)] flex-shrink-0">
            <span className="text-xs text-[var(--text-5)] font-medium uppercase tracking-wider">
              {t.preview}
            </span>
          </div>
          <div className="flex-1 overflow-hidden">
            <DiagramPreview
              imageUrl={imageUrl}
              drawioContent={drawioContent}
              pdfUrl={pdfUrl}
              loading={loading}
              error={error}
            />
          </div>
        </div>
      </main>

      <Tour id="editor" steps={EDITOR_TOUR} />

      {/* ── Status bar ─────────────────────────────────────────────────────── */}
      <footer className="flex items-center justify-between px-4 h-6 border-t border-[var(--border)] flex-shrink-0">
        <span className="text-[10px] text-[var(--text-6)]">
          diagram-as-code · AWS Architecture Diagrams from YAML ·{' '}
          <a
            href="https://fernando.moretes.com"
            target="_blank"
            rel="noopener noreferrer"
            className="text-[var(--text-5)] hover:text-[var(--text-3)] transition-colors"
          >
            {t.footerBy}
          </a>
        </span>
        <span className="text-[10px] text-[var(--text-6)]">
          {t.mode(format)} · <kbd className="font-mono">Ctrl+Enter</kbd>
        </span>
      </footer>
    </div>
  )
}
