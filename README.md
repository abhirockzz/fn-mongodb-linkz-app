# Simple link sharing app using Fn Functions with MongoDB

> similar to stuff like [Hacker News](https://news.ycombinator.com/news

It's possible to

- Create link with details (URL, title and tags)
- Comment and upvote on them
- Search for all content, by user(s) or specifc tag(s)

The app consists of multiple functions corresponding to the above mentioned capabilities and written using a combination of `Java` and `Go`

- `read`, `comment` use the [Fn Java FDK](https://github.com/fnproject/fdk-java) and the [MongoDB Java driver](https://mongodb.github.io/mongo-java-driver/)
- `submit`, `upvote`, `delete` use the [Fn Go FDK](https://github.com/fnproject/fdk-go) and the [MongoDB Golang driver](https://github.com/mongodb/mongo-go-driver)

## Pre-requisites

Clone this repo

### Switch to correct context

- `fn use context <your context name>`
- Check using `fn ls apps`

## Create app

`fn create app --annotation oracle.com/oci/subnetIds=<SUBNETS> --config MONGODB_CONNECT_STRING=<MONGODB_CONNECT_STRING> --config MONGODB_DB=<MONGODB_DB> --config MONGODB_COLLECTION=<MONGODB_COLLECTION> fn-mongodb-app`

e.g.

`fn create app --annotation oracle.com/oci/subnetIds='["ocid1.subnet.oc1.phx.aaaaaaaaghmsma7mpqhqdhbgnby25u2zo4wqlrrcskvu7jg56dryxt3hgvkz"]' --config MONGODB_CONNECT_STRING=mongodb://localhost:27017 --config MONGODB_DB=linksdb --config MONGODB_COLLECTION=links fn-mongodb-app`

**Check**

`fn inspect app fn-mongodb-app`

## Moving on...

Deploy the app...

`cd fn-mongodb-app` and `fn -v deploy --app fn-mongodb-app`

## Test it out..

### Submit a link

`echo -n '{"link":"https://betterexplained.com/articles/intuitive-understanding-of-eulers-formula/","headline":"Intuitive Understanding of Euler’s Formula","postedBy":"foo","tags":["math,science"]}' | fn invoke fn-mongodb-app create`

If successful, you'll get a JSON response

`{"Status":"SUCCESS","Message":"5bd6d40f22baacd8553c2258"}`

`Message` containes the auto generated ID (by MongoDB) of the link you just inserted. You'll use this in the subsequent steps

### Read data


- Filter links by user - `echo -n '{"name":"users", "value":"foo"}' | fn invoke fn-mongodb-app read`

		[
			{
				"link": "https://betterexplained.com/articles/intuitive-understanding-of-eulers-formula/",
				"headline": "Intuitive Understanding of Euler’s Formula",
				"postedBy": "foo",
				"tags": [
					"math",
					"science"
				]
			}
		]

- Filter links by tags -`echo -n '{"name":"tags", "value":"python,science"}' | fn invoke fn-mongodb-app read`

		[
			{
				"link": "https://betterexplained.com/articles/intuitive-understanding-of-eulers-formula/",
				"headline": "Intuitive Understanding of Euler’s Formula",
				"postedBy": "foo",
				"tags": [
					"math",
					"science"
				]
			},
			{
				"link": "https://code.fb.com/ml-applications/qnnpack-open-source-library-for-optimized-mobile-deep-learning/",
				"headline": "Qnnpack: PyTorch-integrated open source library for mobile deep learning",
				"postedBy": "bar",
				"tags": [
					"python",
					"deep learning"
				]
			}
		]

- List all - `echo -n '{"name":"ALL", "value":""}' | fn invoke fn-mongodb-app read`

		[
			{
				"link": "https://betterexplained.com/articles/intuitive-understanding-of-eulers-formula/",
				"headline": "Intuitive Understanding of Euler’s Formula",
				"postedBy": "foo",
				"tags": [
					"math",
					"science"
				]
			},
			{
				"link": "https://code.fb.com/ml-applications/qnnpack-open-source-library-for-optimized-mobile-deep-learning/",
				"headline": "Qnnpack: PyTorch-integrated open source library for mobile deep learning",
				"postedBy": "bar",
				"tags": [
					"python",
					"deep learning"
				]
			}
		]


### Comment

`echo -n '{"linkDocId": "5bd6d40f22baacd8553c2258", "comment":"That's an informative article", "user":"foobar"}' | fn invoke fn-mongodb-app comment`

Success response - `Successfully added comment 5bd6d40f22baacd8553c2258`

### Upvote a link

This will add a vote to the link specified by the ID

`echo -n '5bd6d40f22baacd8553c2258' | fn invoke fn-mongodb-app upvote`

If successful, you'll receive a `Upvoted successfully`  response

### Delete a link

This will delete the link specified by the ID

`echo -n '5bd6d40f22baacd8553c2258' | fn invoke fn-mongodb-app delete`

If successful, you'll receive a `Deleted successfully`  response
