open System
open System.IO

let log x = printfn "%A" x

type Point = int * int
type Move = char * int

let distance ((x, y): Point) = Math.Abs(x) + Math.Abs(y)

let moves (l: string) =
    let move (s: string) = (s[0], int s[1..])
    l.Split "," |> Seq.map move |> Seq.toList

let wire (moves: list<Move>) =
    let follow =
        List.fold (fun ((x, y), w) (d, a) ->
            let p =
                match d with
                | 'L' -> [ for i in x - 1 .. -1 .. x - a -> (i, y) ]
                | 'R' -> [ for i in x + 1 .. x + a -> (i, y) ]
                | 'D' -> [ for i in y - 1 .. -1 .. y - a -> (x, i) ]
                | 'U' -> [ for i in y + 1 .. y + a -> (x, i) ]
                | _ -> failwith "invalid direction"

            (p |> Seq.last, w @ p))

    follow ((0, 0), []) moves |> snd


let input = File.ReadAllLines("in")

let wire1 = moves input[0] |> wire
let wire2 = moves input[1] |> wire

let intersections =
    Set.intersect (set wire1) (set wire2) |> Seq.filter (fun x -> x <> (0, 0))

let ans1 = intersections |> Seq.map distance |> Seq.min
log ans1

let steps x =
    let s1 = List.findIndex (fun y -> y = x) wire1 + 1
    let s2 = List.findIndex (fun y -> y = x) wire2 + 1
    s1 + s2

let ans2 = intersections |> Seq.map steps |> Seq.min
log ans2

