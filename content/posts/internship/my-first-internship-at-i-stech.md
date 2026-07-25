title: What I Learned From My First Software Engineering Internship
date: 2026-07-25
description: Six months at i-stech, two very different projects, and a lot of things I didn't know when I started.
tags: internship, backend, cloud, aws, fastapi, nextjs
---

From January to June 2026, I worked as a Software Engineer Intern at i-stech. It was my first professional experience in software engineering, and I quickly realised that working on real-world software was very different from the projects I had built at university.

I knew how to build APIs, work with databases, and develop web applications through my academic projects. However, I soon learned that these skills represented only a small part of what it takes to work on a real software system.


## PEP

The first project I worked on was PEP, an internal project management platform built with FastAPI, Next.js, Supabase, and Stripe.

At first, the project felt quite familiar. I had built CRUD applications before, so the basic concepts were not new to me. The difference was the amount of business logic involved. There were many cases that I had never really thought about in my university projects.

The payment system was probably the part that taught me the most. PEP needed to support both bank transfers and Stripe payments. Refunds were not always as simple as returning money to a customer. Depending on the situation, the refund could affect the customer's balance and future invoices. There were also Stripe webhooks and background jobs involved in keeping different parts of the system synchronised.

I spent quite a lot of time trying to understand the payment flow before making changes to the code. At the beginning, I felt like I was moving too slowly. I wanted to start implementing things immediately, but I kept finding more cases to understand. Eventually, I realised that spending time understanding the flow was much better than writing a quick implementation and then fixing problems caused by incorrect assumptions.

This was probably the first time I really understood why idempotency matters. An external event might be sent more than once. A background job might fail halfway through. Two operations might happen at almost the same time. These were things I had learned about theoretically, but working with them in an actual system made the concepts much more concrete.

I also learned that background jobs are not something you can simply start and forget about. When something runs automatically in the background, you need to know whether it succeeded, whether it failed, and what happened when it failed. Logging and observability became much more important to me after working on this project.

The biggest lesson from PEP was probably that the difficult part of a feature is not always the code. Sometimes the difficult part is understanding all the possible states of the system.

## TTA Analytics

The second half of my internship was very different. I worked on TTA Analytics, a platform that processed a large amount of socioeconomic data. The system contained 13 relational datasets, with some files containing more than 28 million rows. The data covered demographic and economic forecasts for 1,741 municipalities and 107 industry sectors through 2050.

Before working on this project, my experience with AWS was mostly limited to tutorials and small personal projects. Suddenly, I had to work with a much larger system and think about how the data would be processed and stored in a real environment.

One of the first things I learned was that the approach that works for a small dataset does not necessarily work for a large one. Reading an entire CSV file into memory is straightforward when the file is small. It becomes a problem when the file is hundreds of megabytes or larger.

I had to start thinking more carefully about how data was ingested. Instead of loading everything at once, the data needed to be processed in chunks. Staging tables were used to validate the data before moving it into the final tables. The database schema also needed to be designed with future queries in mind, especially because the system involved joins across multiple large datasets.

This project also gave me a much better understanding of how many different parts are involved in running a system on AWS. Sometimes the problem was not in the application code at all. It was a permission, configuration, or environment issue. I had used AWS before, but working with it as part of a real project was a completely different experience.


There were times when I knew the code was probably correct, but the application still could not access a resource. The problem was somewhere in the configuration. These issues were frustrating, but they forced me to understand more about how the infrastructure actually worked instead of treating AWS as a black box.

## Code Reviews

One of the things I remember most from the beginning of my internship was my first major code review.

I received a lot of comments on the pull request. My first reaction was honestly embarrassment. I had spent a lot of time working on the implementation, so seeing so many things that could be improved was not easy.

After looking through the comments more carefully, though, I realised that most of them were things I could learn from. Some comments were about code quality, some pointed out edge cases I had missed, and others were suggestions that would make the code easier to understand.

Over time, I started to appreciate code reviews more. It is easy to think of feedback as criticism when you are still inexperienced. But having someone with more experience look at your code is also a way to learn things before they become problems in production.

I still don't particularly enjoy opening a pull request and finding a long list of comments. But I definitely prefer receiving those comments during review rather than discovering the same problems after the code has already been deployed.

## Things That Were Not in the Textbooks

One of the biggest differences between university projects and working on real software was everything around the code itself.

At university, if the application worked on my computer, I was usually quite satisfied. During my internship, that was obviously not enough.

I encountered problems that only happened in certain environments. There were differences between local development and staging. There were configuration issues and environment variables that were easy to overlook. At one point, I spent a significant amount of time debugging a problem that was ultimately caused by a difference between the environments.

It was not an especially interesting bug, but it taught me an important lesson: the environment is part of the system too.

I also became much more aware of the importance of communication. When I was stuck, simply saying that something did not work was usually not very helpful. It was much better to explain what I expected to happen, what actually happened, what I had already tried, and where I thought the problem might be.

The same applied to documentation and pull requests. Writing down why I made a particular decision often helped other people understand the change, but it also helped me understand my own thinking more clearly.

## Looking Back

If I could go back to the beginning of my internship, I would probably ask questions earlier. During the first few weeks, I sometimes spent too much time trying to figure out requirements on my own because I was afraid of asking questions that might seem obvious.

In reality, spending hours making assumptions is usually much worse than asking someone a question that takes five minutes to answer.

I would also pay more attention to logs. My first instinct when something went wrong was often to start changing the code immediately. Many times, the logs already contained useful information about what was happening.

I would write more things down as well. During the internship, I made many decisions, tried different approaches, and ran into problems that I later had to think about again. Having a record of those things would have saved me time.

The last thing I would change is the size of my pull requests. Large PRs are difficult to review, difficult to test, and more difficult to debug when something goes wrong. I understood this much better after experiencing it myself.

My internship was not always easy. There were many times when I felt lost, especially when I was working with a system or technology I had never seen before. Looking back, I think the most important thing I learned was that feeling lost is just part of starting something new. There were many times during my internship when I had no idea where to begin. I would spend time reading the existing code, checking the documentation, looking through logs, and trying to understand what was actually happening. Sometimes I found the problem quickly. Sometimes I didn't. Sometimes I made the problem worse before finally understanding it. I even reset our staging database once :)) 

That experience has made me more comfortable with not knowing the answer immediately. I don't expect that feeling to disappear, especially now that I am starting a new role as a Software Engineer Intern at Grab, where I will be working on Mobile development. The technology stack will be different, and I will probably have to learn a lot from the beginning again.

But this time, I know that being unfamiliar with something doesn't necessarily mean I am not capable of doing it. It usually just means I haven't spent enough time with it yet.

And after six months at i-stech, I have a better idea of what to do next when I get stuck.
