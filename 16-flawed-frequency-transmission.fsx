open System
open System.IO

let input = (File.ReadAllLines "in")[0]

let numbers = input |> Seq.map Char.GetNumericValue |> Seq.map int |> Seq.toList

let pattern = [ 0; 1; 0; -1 ]

let onesDigit (n: int) : int = Math.Abs n % 10

let patternForElement (i: int) : seq<int> =
    seq {
        while true do
            for n in pattern do
                for _ in 0..i do
                    yield n
    }
    |> Seq.skip 1

let message (numbers: list<int>) (offset: int) =
    String.Join("", numbers |> Seq.skip offset |> Seq.take 8)

let rec part1 (numbers: list<int>) (phase: int) =
    if phase = 100 then
        numbers
    else
        let numbers' =
            numbers
            |> Seq.mapi (fun i _ ->
                Seq.zip numbers (patternForElement i)
                |> Seq.map (fun (n, p) -> n * p)
                |> Seq.sum
                |> onesDigit)
            |> Seq.toList

        part1 numbers' (phase + 1)

let ans1 = message (part1 numbers 0) 0
printfn $"ans1 {ans1}"
