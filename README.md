# Creating a course platform from scratch with Golang, Stripe, and AWS

![Home Page](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/homepage.png)

## Why This Project Stands Out

Building a course platform isn’t the challenge — picking up a ready-made framework and slapping together some third-party services can get you 80% of the way there. But this project is different.

This is a fully custom-built course platform, crafted from the ground up using Golang, AWS, and Stripe to handle everything from authentication to payments, asset management, and dynamic content generation. Unlike most platforms that lean heavily on third-party authentication, this project includes a modular authentication system that can be reused across different projects with minimal changes. Every piece of this system was built with scalability, security, and real-world production readiness in mind.

If you’re a Go developer or someone looking to see how real-world applications can be built without hiding behind "just use this library," this is the project for you.

## Key Features

### 🚀 Built from Scratch Authentication System

- No third-party auth providers—everything is handled internally.
- Secure user login and signup with encrypted session cookies.
- Password reset with email verification.
- Suspicious activity detection, notifying users of unusual logins.
- Modular enough to be extracted and used in other projects.

### 💰 Full Stripe Integration

- Handle course purchases seamlessly with Stripe.
- Built-in refund and dispute management.
- Webhook handling to update user access in real time.

### 🏗 AWS-Powered Infrastructure

- AWS S3 for storing assets, ensuring efficient and scalable storage.
- AWS CloudFront as a CDN, making content delivery lightning-fast.
- Automatic asset syncing with simple CLI commands.

### ⚡ Blazing Fast Performance (100/100 Lighthouse Score)

- Server-side rendering for near-instant load times.
- Optimized asset delivery through CloudFront.
- Efficient database queries to minimize response times.

### 📡 Full RSS Feed Support

- Users can subscribe to an RSS feed for new tutorials and courses.
- Updates are generated dynamically based on available content.

### 📜 Auto-Generated Sitemaps

- The platform dynamically generates a full sitemap at runtime using [sitemapper](https://github.com/PsionicAlch/sitemapper).
- Ensures that all public URLs are indexed correctly for SEO.

### 🎓 User Certification System

- Users receive a PDF certificate upon course completion.
- Certificates are generated client-side, reducing server load.

### 🛠 Admin Dashboard

- Fully functional admin panel to manage all aspects of the platform.
- Control comments, courses, discounts, purchases, refunds, tutorials, and users from one place.
- Designed for scalability and ease of use.

### 🌟 Affiliate Program

- Users can share their affiliate code to earn discounts on future course purchases.
- Incentivizes engagement and word-of-mouth marketing.
- Fully integrated with Stripe for seamless discount handling and price calculations.

### 🔧 Comprehensive User Settings Page

- Users can update their first and last names.
- Change their email address and password.
- Manage their whitelisted IP addresses for added security.
- Request refunds directly from their settings.
- Permanently delete their profile if needed.

This project isn't just about writing code—it's about engineering a real-world solution. Whether you're looking to learn Golang, explore scalable system design, or just see what a true end-to-end application looks like without cutting corners, this project delivers.

![Administration Page](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/adminpage.png)

![Login Page](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/loginpage.png)

![Lighthouse Performance](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/performancepage.png)

![Course Purchase Page](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/purchasepage.png)

![RSS Feed](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/rsspage.png)

![User Settings Page](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/settingspage.png)

![Sitemap](https://github.com/PsionicAlch/course-platform/blob/main/screenshots/sitemappage.png)

## How to get this project running on your local device?

### Before you begin:

Before you start, you will need to set up a few external things. You will need the following:

- Stripe account
- Stripe CLI
- AWS Access Key
- AWS S3 Bucket
- AWS CloudFront

The **Stripe account** is used to simulate payments with the secret key. A test account is enough for local development purposes.

The **Stripe CLI** is required to simulate payments on your local system.

**AWS S3** and **CloudFront** are used as a CDN to serve the assets for this project. Setting both up is completely free and I will point you to the following tutorial: https://aws.amazon.com/cloudfront/getting-started/S3/

The **AWS Access Key** is required to sync the local assets with your S3 bucket.

### Step 1: Create a .env file

All the configuration for this project lives in a .env file. You can copy the following one and just tweak it as necessary:

```env
PORT=8080
ENVIRONMENT=development
DOMAIN_NAME=localhost

NOTIFICATION_COOKIE_NAME=notifs
AUTH_COOKIE_NAME=auth
AUTH_TOKEN_LIFETIME=43200
EMAIL_TOKEN_LIFETIME=30
CURRENT_SECURE_COOKIE_KEY=
PREVIOUS_SECURE_COOKIE_KEY=

EMAIL_PROVIDER=smtp
EMAIL_HOST=localhost
EMAIL_PORT=1025
EMAIL_ADDRESS=contact@example
EMAIL_PASSWORD=

STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=

CLOUDFRONT_URL=
REGION=
ACCESS_KEY_ID=
SECRET_ACCESS_KEY=
BUCKET_NAME=
```

**PORT**: The port that the server will run on. 

**ENVIRONMENT**: The current environment of the project. Only "development", "testing", and "production" are viable options for this. Anything else will result in an error.

**DOMAIN_NAME**: The domain name of the project.

**NOTIFICATION_COOKIE_NAME**: The name you want the notification cookie to have. This cookie is only used for flash messages and nothing else.

**AUTH_COOKIE_NAME**: The name you want the authentication cookie to have.

**AUTH_TOKEN_LIFETIME**: How long you want an authentication token to be valid for. This number is in minutes: 60 minutes per hour * 24 hours per day * 30 days = 43200.

**EMAIL_TOKEN_LIFETIME**: How long you want an email token to be valid for (email tokens are used when resetting a user's password): This number is in minutes.

**CURRENT_SECURE_COOKIE_KEY**: This is a string used to encrypt the authentication cookie's contents. This project makes use of the [Gorilla Secure Cookie](github.com/gorilla/securecookie) package. A new key can be generated by typing ```make new-key``` in your terminal.

**PREVIOUS_SECURE_COOKIE_KEY**: This is the previous key you used. It's recommended to swap out your keys on a somewhat regular basis. NOTE: It can be left empty.

**EMAIL_PROVIDER**: The email provider you want to use. Currently "smtp" is the only valid option.

**EMAIL_HOST**: The host URL of your SMTP server. I used [MailHog](https://github.com/mailhog/MailHog) during testing.

**EMAIL_PORT**: The port of your SMTP server.

**EMAIL_ADDRESS** The email address that your emails will come from.

**EMAIL_PASSWORD**: The password to authenticate your email with. NOTE: It can be left empty.

**STRIPE_SECRET_KEY**: Your Stripe secret key. You can find this in your Stripe dashboard. The test key is enough.

**STRIPE_WEBHOOK_SECRET**: Your Stripe webhook secret key. This can be found in your Stripe dashboard. If you are just working locally Stripe CLI will give you one to use.

**CLOUDFRONT_URL**: The URL for your AWS CloudFront instance.

**REGION**: The region where your AWS S3 bucket is currently hosted (eg "eu-west-3").

**ACCESS_KEY_ID**: The access key ID for your AWS.

**SECRET_ACCESS_KEY**: The secret key that came with your AWS access key.

**BUCKET_NAME**: The name of your AWS S3 bucket.

### Step 2: Syncing local assets with AWS

You will need to sync the local assets to your AWS S3 bucket so that all the assets are visible on your side. This can be done with a simple command in your terminal:

```bash
make sync-assets
```

### Step 3: Migrating your database

You will need to set up a local copy of the SQLite database. This can be done with a simple command in your terminal:

```bash
make migrate-up
```

### Step 4: Loading all tutorials and courses into your database

When creating this project it was important for me to always have a local copy of all the tutorials and courses so that I could save them to git. All tutotrials and courses are written in Markdown and FrontMatter so they will first need to be parsed and synced in the database. This can be done with a simple command in your terminal (it's also pretty fast which is something I'm quite proud of):

```bash
make load-content
```

### Step 5 (optional): Seeding the database with some dummy content

To play around with the project in local development I created some seed scripts to add dummy users and discounts to the database. This can be done with a simple command in your terminal:

```bash
make seed-database
```

### Step 6: Running the project locally

To help speed up local development I used [air](https://github.com/air-verse/air) for live reloading. To do this on your side as well you will need to have [air](https://github.com/air-verse/air) installed on your local development system. This can be done with a simple command in your terminal:

```bash
go install github.com/air-verse/air@latest
```

Once you have [air](https://github.com/air-verse/air) installed you start it up with the following command in your terminal:

```bash
air
```

If you don't want to use [air](https://github.com/air-verse/air) you can also just run the application locally by using the following command in your terminal:

```bash
make run
```

This will build the project and run it.

## How to write a tutorial?

To write a tutorial you will need to create a Markdown file under ./web/content/tutorials. The name you give your Markdown file does not matter at all so you can use whatever works best for. This project uses FrontMatter for setting file based metadata so your file should always start with something like this:

```markdown
---
title: "A SEO FRIENDLY TITLE FOR YOUR TUTORIAL"
description: "A SEO FRIENDLY SHORT DESCRIPTION OF YOUR TUTORIAL"
thumbnail_url: "THE EXACT URL PATH FOR YOUR TUTORIAL'S THUMBNAIL IMAGE"
banner_url: "THE EXACT URL PATH FOR YOUR TUTORIAL'S BANNER IMAGE"
keywords: ["THE LIST OF SEO KEYWORDS YOU WANT YOUR TUTORIAL TO HAVE"]
key: "A UNIQUE STRING TO ASSOCIATE WITH THIS SPECIFIC TUTORIAL"
---
```

You can generate a key using the following command: ```make generate-file-key```. Each file key should be unique because it is used in the database to uniquely identify each tutorial. If two or more tutorials share the same file key they will override each other in the database.

The rest of the tutorial can be written in plain Markdown. Images are supported but you will need to provide the exact URL path to the image. The reason for this is so that you can use images that aren't hosted by you. You can also write code blocks and they will be properly sytnax highlighted using [highligh.js](https://highlightjs.org/).

Once your tutorial has been written you can load it into the database with the following command: ```make load-content```. The tutorial will be set to "unpublished" by default without an author so you will need to [publish your tutorial]("https://github.com/PsionicAlch/course-platform?tab=readme-ov-file#how-to-publish-a-tutorial") before it's visible.

## How to publish a tutorial?

To publish a your newly created tutorial you will need have an admin account. You can create a new admin account with the following command: ```make new-admin name="YOUR NAME" surname="YOUR SURNAME" email="YOUR EMAIL ADDRESS" password="YOUR PASSWORD"```. After the new admin account has been registered you'll need to login using the email and password you just set.

Now that you're logged into your admin account you'll notice a new "Admin" link in your navbar. You can click on it to go to the admin panel. You'll be redirected to "/admin/users". This is the users administration panel. You'll need to set at least one user to be an author. You can do this by scrolling to the column called "Author". This should currently be a cross to indicate that the user you're looking at is not an author. You can double click the cell to get a dropdown. Select "Author" from the dropdown. If all goes well the user's status should be set to "Author" and you should now see a checkmark.

Next you need to go to the tutorials administration panel. You can do this by clicking on the hamburger icon on the right of your screen. A menu should popup. You will need to click on "Tutorial Managment" to be redirected to the tutorials administration panel. Find the tutorial you want to publish in the table and scroll to the "Published" column. You can double click the cell to get another popup. Select "Published" from the dropdown menu. Lastly scroll to the "Author" column and double click the cell. You should get another dropdown menu. Select an author from the dropdown menu. 

Now that your tutorial has been set to "Published" and have an author you can go to "/tutorials" and you should see your tutorial there.

## How to write a course?

To write your first course you can head over to the ./web/content/courses folder of the source code. 

A course is made up of two parts. The first part is a Markdown file that represents the sales page for your course. The second part is a folder that contains all the Markdown files. Each Markdown files represents an individual chapter of your course. The name of the sales page Markdown file does not matter, nor does the folder that contains the chapters so you can name it whatever works best for you.

The sales page file uses FrontMatter for the file specific metadata. So the start of your file should always contain the following metadata:

```markdown
---
title: "A SEO FRIENDLY TITLE OF THE COURSE"
description: "A SEO FRIENDLY SHORT DESCRIPTION OF THE COURSE"
thumbnail_url: "THE EXACT URL PATH FOR THE COURSE'S THUMBNAIL IMAGE"
banner_url: "THE EXACT URL PATH FOR THE COURSE'S BANNER IMAGE"
keywords: ["THE LIST OF SEO KEYWORDS YOU WANT YOUR COURSE TO HAVE"]
key: "A UNIQUE STRING TO ASSOCIATE WITH THIS SPECIFIC COURSE"
---
```

You can generate a key using the following command: ```make generate-file-key```. Each file key should be unique because it is used in the database to uniquely identify each course and to link the chapter files that are associated with this course. If two or more courses share the same file key they will override each other in the database.

The rest of the file should be used to write an effective sales page for the course. Don't include information like the price of the course nor any links to the buy button since that will automatically be added.

Next up this the chapter files. The name of folder containing the individual chapters doesn't matter. Each chapter files uses FrontMatter for the file specific metadata. So the start of each chapter file should always contain the following metadata:

```markdown
---
title: "THE TITLE OF THIS SPECIFIC CHAPTER"
chapter: "AN INTEGER THAT STATES THE CURRENT CHAPTER NUMBER"
course_key: "THE FILE KEY OF THE COURSE SALES PAGE FILE THAT THIS CHAPTER IS LINKED TO"
key: "A UNIQUE STRING TO ASSOCIATE WITH THIS SPECIFIC CHAPTER"
---
```

You can generate a key using the following command: ```make generate-file-key```. Each file key should be unique because it is used in the database to uniquely identify each chapter. If two or more chapters share the same file key they will override each other in the database.

This file can contain images and code blocks and the code blocks will be properly sytnax highlighted using [highligh.js](https://highlightjs.org/).

Once your course has been written you can load it into the database with the following command: ```make load-content```. The course will be set to "unpublished" by default without an author so you will need to [publish your course]("https://github.com/PsionicAlch/course-platform?tab=readme-ov-file#how-to-publish-a-course") before it's visible.

## How to publish a course?

To publish a your newly created course you will need have an admin account. You can create a new admin account with the following command: ```make new-admin name="YOUR NAME" surname="YOUR SURNAME" email="YOUR EMAIL ADDRESS" password="YOUR PASSWORD"```. After the new admin account has been registered you'll need to login using the email and password you just set.

Now that you're logged into your admin account you'll notice a new "Admin" link in your navbar. You can click on it to go to the admin panel. You'll be redirected to "/admin/users". This is the users administration panel. You'll need to set at least one user to be an author. You can do this by scrolling to the column called "Author". This should currently be a cross to indicate that the user you're looking at is not an author. You can double click the cell to get a dropdown. Select "Author" from the dropdown. If all goes well the user's status should be set to "Author" and you should now see a checkmark.

Next you need to go to the courses administration panel. You can do this by clicking on the hamburger icon on the right of your screen. A menu should popup. You will need to click on "Course Managment" to be redirected to the courses administration panel. Find the course you want to publish in the table and scroll to the "Published" column. You can double click the cell to get another popup. Select "Published" from the dropdown menu. Lastly scroll to the "Author" column and double click the cell. You should get another dropdown menu. Select an author from the dropdown menu. 

Now that your course has been set to "Published" and have an author you can go to "/courses" and you should see your tutorial there.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/PsionicAlch/course-platform/blob/main/LICENSE) file for details.