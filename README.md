# Dad :man_shrugging:
[![Go Report Card](https://goreportcard.com/badge/github.com/alee792/dad)](https://goreportcard.com/report/github.com/alee792/dad) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-89%25-brightgreen.svg?longCache=true&style=flat)</a>
### Dad is good at retelling jokes, but he's not so good at coming up with his own.

## Try it out here!
["Hey, I heard a good joke the other day..."](https://dad-git-master.alee792.now.sh/joke)

["Is this the kind of stuff you do at work?"](https://dad-git-master.alee792.now.sh/hn)

*Examples are 1grams*

Dad reads jokes on icanhazdadjoke.com, and then tries to come up with his own. Results can be...interesting.
> One about clocks. It's because I woke up the sea lion?  
> The hokey pokey, but there is really need to soap, but now that’s a forest and put down!
> Nailed it!  

Sometimes, Dad will see a cool story on [HackerNews](https://news.ycombinator.com/news) and try to talk to you about it. Oddly enough, these mangled headlines sound like pitches from Y Combinator's most recent cohort.
> To bring deep learning to use Outlook  
> To a month, Without NYT or Go?  
> 11 Press Kits from Python APIs with specially crafted D-Bus message  
> Salt water tanks pose health data to be taught code  
> Powered by PHP (2017)  

## Thanks 
Inspired by this [codewalk](https://golang.org/doc/codewalk/markov/). Enhanced with configurable n-gram levels and a more legible/less clever implementattion.

And [gomarkov](https://github.com/mb-14/gomarkov) for the idea of spoofing HackerNews headlines.

## Usage
### Compile
If you have Go, `make bin` or `go build ./cmd/dad`. This project does use Modules, so build accordingly!
If you have a Mac, you can use the `dad` binary in the project root.

## Run
`./dad` runs an HTTP server with routes of `/joke` and `/realjoke`, where you can `GET`, respectively, butchered Markov jokes or jokes directly from the corpus' source.

To run in HackerNews mode, simply supply a flag, `./dad -s=hn`. The routes are the same, but can be moved in the future.

By default, Dad runs using 2grams. Set n to change the order of the ngrams, e.g. `./dad -n=1`

## Caching
N-grams are automatically saved and loaded in `./bin` with the format `{n}grams-{source}.json`. Some prepopulated 1 and 2 grams are included.

## Future work
1. More memory considerate data structure.
2. Additional flags for port, warm up size, etc.
3. Allow query params to dictate ngram order.
