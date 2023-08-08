open System
open System.IO
open System.Text.Json
open System.Text.RegularExpressions

let log x =
    let options = JsonSerializerOptions(WriteIndented = true)
    let json = JsonSerializer.Serialize(x, options)
    printfn "%s" json

type Reagent = { Quantity: int64; Chemical: string }
type Reaction = { Left: Set<Reagent>; Right: Reagent }

type Reactions = Map<string, Reaction>
type Requirements = Map<string, int64>

let parseReaction (s: string) : Reaction =
    let reagent (m: Match) : Reagent =
        { Quantity = int m.Groups[1].Value
          Chemical = m.Groups[2].Value }

    let pattern = "(\d+) ([A-Z]+)"
    let parts = s.Split "=>"

    let left = Regex.Matches(parts[0], pattern) |> Seq.map reagent |> Set.ofSeq
    let right = reagent (Regex.Match(parts[1], pattern))

    { Left = left; Right = right }

let makeReactions (lines: seq<string>) : Reactions =
    lines
    |> Seq.map parseReaction
    |> Seq.map (fun r -> r.Right.Chemical, r)
    |> Map.ofSeq

let nonOre (requirement: Requirements) : Requirements =
    requirement |> Map.filter (fun k _ -> k <> "ORE")

let nonZero (requirement: Requirements) : Requirements =
    requirement |> Map.filter (fun _ v -> v > 0)

let addRequirement (n: int64) (requirements: Requirements) (reagent: Reagent) : Requirements =
    let amt = n * reagent.Quantity

    requirements
    |> Map.change reagent.Chemical (fun x ->
        match x with
        | Some qty -> Some(qty + amt)
        | None -> Some(amt))

let rec produce (reactions: Reactions) (requirements: Requirements) : Requirements =
    let outstanding = requirements |> nonZero |> nonOre

    if outstanding.IsEmpty then
        requirements
    else
        let chemical, qty = (outstanding |> Seq.minBy (fun kvp -> kvp.Key)).Deconstruct()
        let reaction = reactions[chemical]
        let productQty = reaction.Right.Quantity
        let applications = max 1L (qty / productQty)

        let qty' = qty - (applications * productQty)
        let required' = requirements |> Map.add chemical qty'
        let required'' = Set.fold (addRequirement applications) required' reaction.Left

        produce reactions required''

let oreForFuel (reactions: Reactions) (n: int64) : int64 =
    let required0 = Map [ ("FUEL", n) ]
    let required = produce reactions required0
    required["ORE"]

let part1 reactions = oreForFuel reactions 1

let reactions = makeReactions (File.ReadAllLines "in")

let ans1 = part1 reactions
printfn $"ans1: {ans1}"

let oreLimit = 1000000000000L

let rec searchFuel (reactions: Reactions) (lower: int64) (upper: int64) =
    let mid = (upper + lower + 1L) / 2L
    let ore = oreForFuel reactions mid

    if upper = lower then
        upper
    else if ore > oreLimit then
        searchFuel reactions lower (mid - 1L)
    else
        searchFuel reactions mid upper

let rec findUpper reactions n =
    let ore = oreForFuel reactions n
    if ore > oreLimit then n else findUpper reactions (n * 2L)

let part2 reactions =
    let base' = oreForFuel reactions 1
    let upper = findUpper reactions (oreLimit / base')
    searchFuel reactions (upper / 2L) upper

let ans2 = part2 reactions
printfn $"ans2: {ans2}"
