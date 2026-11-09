# RFC 17 — Link expiry lab

**Covered RFC 17. Next RFC 18.** RFC 16's public main branch was
[`a6b877d`](https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0016-distribution-lab/commit/a6b877d4b29cbdfed8279119b9bf4494ed9a0c21)
when checked. Its README, command, module, license and workflow were reviewed
before reuse. The RFC Editor lists number 14 as Not Issued.

## Official document

On 27 August 1969, J.E. Kreznar asked what would happen when the early
ARPANET's Interface Message Processor (IMP) dropped idle link-table entries,
acknowledged messages before delivery to a host, or discarded startup
messages. The RFC includes R.E. Kahn's reply as RFC 17a. Kahn explained that
an idle entry can be recreated when the table has space and that blocked
transmit entries survive expiry until the Ready for Next Message (RFNM)
signal releases them. He placed RFNM at the head of the destination IMP's
host-output queue, just before transfer to the host, and explained that
the startup NOP was a disposable control message.

The printed header says **Informational**. The RFC Editor machine-readable
catalog lists the **Legacy** stream and **Unknown** publication and current
status, with empty updates, updated-by, obsoletes and obsoleted-by lists.
Checked 9 November 2026. This exchange documents early ARPANET design.

- [Official text and reply](https://www.rfc-editor.org/info/rfc17/)
- [RFC Editor catalog metadata](https://www.rfc-editor.org/rfc/rfc17.json)
- [RFC Editor index](https://www.rfc-editor.org/rfc-index.txt)

The modern lesson is to specify what a timer or acknowledgment actually
means. An idle slot can be reclaimed without implying that the logical link
has disappeared. A signal generated at a queue boundary can precede delivery.

## Applications

1. **Connection pool simulator, selected.** Show how idle expiry frees a
   slot while an in-flight operation retains one during a long pause.
2. **Banking API event trace.** Give acceptance by an intermediary and
   completion by the downstream worker separate events and timestamps.
3. **Cloud gateway capacity exercise.** Replay competing clients against a
   bounded table and distinguish contention on one slot from a full table.

## Run locally

Use Go 1.22 or later. From the repository root, run the following commands.

```sh
go run ./cmd/linkexpiry < examples/saturation.json
go run ./cmd/linkexpiry -html < examples/saturation.json > docs/index.html
go test -race ./...
go vet ./...
go build ./...
```

In the sample, link `b` expires at 32 seconds. Link `a` remains active and
link `c` takes the open slot. At 33 seconds, trying to recreate `b` is
delayed by the full table. Link `a` is released at 40 seconds and its idle
entry expires at 70 seconds, when `b` can be inserted again.

## Model

This is a deterministic teaching simulation of one table behavior described
in RFC 17a. The script models unblocking as a synthetic `release` event
and expiry as a sweep before each input event. The CLI reads a JSON object
with capacity 1–64 and 1–500 ordered events. Each event has a nonnegative
`at_seconds`, an action of `begin` or `release`, and a short lowercase link
identifier. Events at the same time run in input order. A delayed begin
leaves the table state alone after expiry and can be retried by adding a
later event. The input limit is 1 MiB.

Logical time and bounded inputs make experiments repeatable. Real clock
scheduling, queue delivery, RFNM encoding and host/IMP wire behavior are
outside the simulation. GitHub Pages serves a saved HTML trace generated
by the local Go CLI.

## Validation

The suite checks idle expiry at the 30-second boundary, protection for an
active transfer, saturation, reinsertion, reuse, invalid scripts and JSON
and HTML output. CI also checks that regenerating the published HTML produces
the committed page. The source is MIT-licensed.
