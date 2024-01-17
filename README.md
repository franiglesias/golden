# Golden

A Go library for snapshot 📸 testing.

## TL;DR

> Current Status: not ready for production

### Installation

No, seriously. It's not ready ✋.

```shell
not yet
```

### Basic Usage:

```go
func TestSomething(t *testing.T) {
	output := SomeFunction("param1", "param2")
	
	golden.Verify(t, output)
}
```

## What is Golden?

Golden is a library inspired by projects like [Approval Testing](https://approvaltests.com/). There are some other similar libraries out there, such as [Approval Tests](https://github.com/approvals/go-approval-tests), [Go-snaps](https://github.com/gkampitakis/go-snaps) or [Cupaloy](https://github.com/bradleyjkemp/cupaloy).

So... Why to reinvent the wheel?

First of all, why not? I was willing to start a little side project to learn and practice some Golang things that I didn't find opportunity during the daily work. For example, creating a library for distribution, some questions about managing state, creating fluent APIs, managing unknown types, etc. 

Second. I found some limitations in the libraries I was using (Approval tests, mainly) that make the work a bit uncomfortable. So, I started to look for alternatives. I wanted some more flexibility and customization.  

## Snapshot testing

Snapshot testing is a testing technique that provides an alternative to assertion testing. In assertion testing, you compare the output of executing some unit of code with some expectation you have about it. For example, you could expect that the output would equal some value, or that it contains some text, etc.

```go
assert.Equal(t, "Expected", output)
```

This works very well for TDD and testing simple outputs. But it is tedious if you need to test complex objects, generated files and other big outputs. Also, it is not a good tool for testing code that you don't know well or that was not created with testing in mind.

In snapshot testing, instead, you first obtain and persist the output of the execution. This is what we call a snapshot.

Then, you make some changes in the affected code, execute the unit again and, finally, you compare this output with the one you persisted. As you can guess, in order to make the test pass, the code changes should not affect the behavior. This way, it is easy to understand that snapshot testing is great to put existing working code under test.

But testing legacy or unknown code is not the only use case for snapshot testing.

In fact, snapshot testing is perfect for testing complex objects or outputs, such as html, xml, json, generated code, etc. Provided that you can create a file with a serialized representation of the response, you can compare the snapshot with subsequent executions of the same unit of code. I mean: suppose you need to generate an API response with lots of data. Instead of trying to figure out how to check every field value, you generate a snapshot with the expected data. After that, you will be using that snapshot to verify the execution. 

But this way of testing is weird for developing new features... How can I know that my generated output is correct so far?

### Approval testing

Approval testing is a variation of snapshot testing. Usually, in snapshot testing, the first time you execute the test, a new snapshot is created and the test automatically passes. Then, you make changes to the code and use the snapshot to ensure there are not behavioral changes. By the way, this is the default behavior of `Golden`.

In Approval Testing, the first test execution is not automatically passed. Instead, the snapshot is created, but you should review and approve it explicitly. This step allows you to be sure that the output is what you want it to be. You can show it to the business people, to your client, to the user of your API or to whoever can review it. Once you get the approval of the snapshot and re-run the test, it will pass.

In fact, you could make changes and re-run the test until you are satisfied with the output and get the approval.

I think that Approval Testing was first introduced by [Llewellyn Falco](https://twitter.com/llewellynfalco). You can [learn more about this technique in their website](https://approvaltests.com/), where you can find how to approach development with it.

Approval testing is a planned feature for Golden.

### Golden master

There is another variation of snapshot testing. **Golden Master** is a technique introduced by Michael Feathers for working with legacy code that you don't understand. With this technique you could achieve 100% coverage really fast, so you can be sure that refactoring will be safe because you always will know if behaviour of the code is broken due to a change you introduced. And the best thing is that you don't need to really understand the code. Once you start refactoring thins, it will be easier to introduce classic assertion testing and probably remove the Golden Master tests.

It consists in the creation of a lot of tests for the same unit of code, introducing combinations of the parameters that you need to pass to such unit. The original Approval Tests library includes Combinations, a library that helps you to generate those combinatorial tests. There are several techniques to guess the best values you can use. You can study the code and search for values in conditionals, for example, with the help of a graphical coverage tool that shows you what parts of the code are executed or not depending on the values.

Once you complete the collection of possible values for each parameter, you will use the combination tools and a lot of tests will be generated for you. The amount of tests is the product of multiplying the number of values per parameter. You can easily achieve tenths and even hundreds of tests for the same code unit.

As you have probably guessed, Golden takes its name from this technique... and because it starts with "Go". (Anyway, I've just found that another golden package exists, so it is possible that I need to change the name.) 

Combinatorial testing is a planned feature for Golden.

## Problems with snapshot testing

Non-deterministic output: this is not an exclusive problem of snapshot testing. Managing non-deterministic output is always a problem. In assertion testing you can introduce property based testing: instead of looking for exact values, you can look for desired properties of the output.

In snapshot testing things are a bit more complicated. It is difficult to check properties of a specific part of the output and ignore the value. Anyway, one solution is to look for specific patterns and do something about them: replace with a fixed but representative value, replace with some reminder... maybe it is possible to ignore that part of the output in order to compare with the snapshot.

Replacement of non-deterministic data is a planned feature for Golden.
