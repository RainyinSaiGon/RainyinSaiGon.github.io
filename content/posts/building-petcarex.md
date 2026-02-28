title: Building PetCareX: What I Learned Deploying a Go API
date: 2026-01-05
description: Lessons from architecting, building, and deploying a pet-care management API with Go, PostgreSQL, and AWS ECS.
tags: Go, AWS, Backend
---

<h2>What is PetCareX?</h2>

<p>PetCareX is a pet health management platform that lets owners track vaccinations, vet visits, and health alerts for their pets. The backend is a REST API written in Go, backed by PostgreSQL, deployed on AWS ECS Fargate.</p>

<h2>Architecture overview</h2>

<ul>
<li><strong>API layer:</strong> Go + <code>chi</code> router, standard <code>net/http</code> server</li>
<li><strong>Database:</strong> PostgreSQL on RDS (Multi-AZ), migrations managed with <code>golang-migrate</code></li>
<li><strong>Auth:</strong> JWT issued by the Go service; short-lived access tokens + refresh tokens stored in Redis</li>
<li><strong>Deployment:</strong> Docker image → ECR → ECS Fargate behind an ALB</li>
<li><strong>Observability:</strong> structured JSON logs (zerolog), custom metrics to CloudWatch</li>
</ul>

<h2>Decision: stdlib vs framework</h2>

<p>I chose <code>chi</code> over Gin or Echo because it composes cleanly with standard <code>http.Handler</code> middleware. If you write your middleware for <code>http.Handler</code>, you can swap routers later without rewriting business logic. Gin's context type breaks this — middleware written for Gin is Gin-specific.</p>

<h2>Database connection pooling gotcha</h2>

<p>On RDS, the default max connections for a <code>db.t3.micro</code> is ~87. With multiple ECS tasks each having a pool of 25, I hit connection limits under load. The fix was setting sensible pool maximums per task:</p>

<pre><code>db.SetMaxOpenConns(10)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
</code></pre>

<p>For anything beyond a few tasks, PgBouncer or RDS Proxy becomes necessary.</p>

<h2>Graceful shutdown</h2>

<p>ECS sends a SIGTERM before killing a container. Go's <code>http.Server.Shutdown</code> drains in-flight requests cleanly:</p>

<pre><code>quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
&lt;-quit

ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
srv.Shutdown(ctx)
</code></pre>

<h2>SHAP explainability layer</h2>

<p>The most interesting part: a Python sidecar runs SHAP analysis on top of an RF model trained on pet health records, exposing a small FastAPI endpoint. The Go service calls it asynchronously and attaches feature importance to vet visit summaries. This lets vets see <em>why</em> the model flagged a pet as high-risk — which was the whole point of the XAI component.</p>

<h2>What I'd do differently</h2>

<ul>
<li>Define an internal <code>errors</code> package from day one rather than sprinkling <code>fmt.Errorf</code> everywhere. Structured errors with codes make API responses much cleaner.</li>
<li>Use database-level enums for status fields instead of stringly-typed constants. Caught a typo bug in staging that would have been a compile error.</li>
<li>Write integration tests against a real database (testcontainers-go) rather than mocking the DB layer.</li>
</ul>
