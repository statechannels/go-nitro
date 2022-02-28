**Title** ☝️ Clear, single line description of the pull request, written as though it was an order.

**Body** A summary of the change and which issue is fixed, including

- what is the problem being solved
- why this is a good approach
- what are the shortcomings of this approach, if any

It may include

- background information, such as "Fixes #"
- benchmark results
- links to design documents
- dependencies required for this change

The summary may omit some of these details if this PR fixes a single ticket that includes these details.

**Changes** [Optional]

Multiple changes are not recommended, so this section should normally be omitted. Unfortunately, they are sometimes unavoidable. If there are multiple logical changes, list them separately.

1. `'foo'` is replaced by `'bar'`
2. `'fizzbuzz'` is optimized for gas

**How Has This Been Tested?** [Optional]

Did you need to run manual tests to verify this change? If so, please describe the tests that you ran to verify your changes. Provide instructions so we can reproduce. Please also list any relevant details for your test configuration

**⚠️ Does this require multiple approvals?** [Optional]

Please explain which reason, if any, why this requires more than one approval.

- [ ] Is it security related?
- [ ] Is it a significant process change?
- [ ] Is it a significant change to architectural, design?

---

#### Code quality

- [ ] I have written clear commit messages
- [ ] I have performed a self-review of my own code
- [ ] This change does not have an unduly wide scope
- [ ] I have separated logic changes from refactor changes (formatting, renames, etc.)
- [ ] I have commented my code wherever necessary (can be 0)
- [ ] I have added tests that prove my fix is effective or that my feature works, if necessary
- [ ] I have added comprehensive [godoc](https://go.dev/blog/godoc) comments 

#### Project management

- [ ] I have assigned myself to this PR
- [ ] I have linked the appropriate github issue
- [ ] I have assigned this PR to the appropriate GitHub project
- [ ] I have assigned this PR to the appropriate GitHub Milestone
