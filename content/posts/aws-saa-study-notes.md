title: AWS SAA-C03 Study Notes
date: 2026-02-10
description: Key concepts and mnemonics I used to pass the AWS Solutions Architect Associate exam on my first attempt.
tags: AWS, Cloud
---

<h2>Why I took the AWS SAA-C03</h2>

<p>After deploying PetCareX on EC2 and ECS I realised I was making infrastructure decisions mostly by instinct. The SAA certification forced me to understand <em>why</em> you pick one service over another, not just how to use it.</p>

<h2>Core domains and what actually matters</h2>

<h3>1. Design Resilient Architectures (26%)</h3>
<ul>
<li><strong>Multi-AZ vs Multi-Region:</strong> Multi-AZ is HA (automatic failover, same data). Multi-Region is DR (manual or routed failover, possible replication lag).</li>
<li><strong>RDS Multi-AZ:</strong> synchronous standby — zero data loss. Read Replicas are async — use them for read scaling, not HA.</li>
<li><strong>ELB health checks:</strong> ALB checks at the target group level; NLB passes TCP through. Know when each is appropriate.</li>
</ul>

<h3>2. Design High-Performing Architectures (24%)</h3>
<ul>
<li><strong>S3 prefix trick:</strong> randomise key prefixes to spread request load across partitions and avoid the 5,500 GET/3,500 PUT per-prefix limit.</li>
<li><strong>ElastiCache vs DAX:</strong> DAX is purpose-built for DynamoDB microsecond reads. ElastiCache (Redis/Memcached) is a general-purpose cache tier.</li>
<li><strong>EFS vs EBS:</strong> EFS mounts to many EC2s (shared NFS); EBS attaches to one. EFS is more expensive per GB but scales automatically.</li>
</ul>

<h3>3. Design Secure Applications (30%)</h3>
<ul>
<li><strong>KMS vs SSM Parameter Store vs Secrets Manager:</strong> KMS = key management. SSM Parameter Store = free config/secrets storage. Secrets Manager = auto-rotation, costs $0.40/secret/month.</li>
<li><strong>SCP vs IAM policies:</strong> SCPs set the outer boundary for the entire AWS account (even root). IAM policies grant permission within that boundary.</li>
<li><strong>VPC endpoints:</strong> Gateway endpoints (S3, DynamoDB) are free. Interface endpoints cost ~$0.01/hr each. Use them to keep traffic off the public internet.</li>
</ul>

<h3>4. Design Cost-Optimised Architectures (20%)</h3>
<ul>
<li>Reserved Instances vs Savings Plans: Savings Plans are more flexible (apply to any EC2 family/region) but also less discount than Standard RIs.</li>
<li>Spot Instances: up to 90% cheaper but can be interrupted. Use for stateless, fault-tolerant workloads — batch jobs, ML training, rendering.</li>
<li>S3 Intelligent-Tiering for unpredictable access patterns; S3 Glacier Instant Retrieval for archives you access a few times a year.</li>
</ul>

<h2>Mnemonics that actually helped</h2>

<p><strong>PERM</strong> for EC2 instance store vs EBS: <em>Persistent = EBS, Ephemeral = instance store, Restartable (survives stop/start) = EBS, Machine-attached (physically) = instance store</em>.</p>
<p><strong>The 3 Ns of NACLs:</strong> <em>Numbered rules, No state (stateless), Network-level</em> — versus Security Groups which are stateful and instance-level.</p>

<h2>Resources I actually used</h2>
<ul>
<li>Stephane Maarek's Udemy course — the best paced resource, good for first pass</li>
<li>Tutorials Dojo practice exams — harder than the real exam, great for identifying gaps</li>
<li>AWS FAQ pages for S3, EC2, RDS — dry but authoritative</li>
</ul>

<h2>Result</h2>
<p>Passed first attempt with 812/1000. The hardest questions were around hybrid connectivity (Direct Connect + VPN) and cost optimisation trade-offs. Study those thoroughly if you're not coming from a networking background.</p>
