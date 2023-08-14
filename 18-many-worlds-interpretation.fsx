open System
open System.IO

type Point = int * int
type Grid = Map<Point, char>

let isEntrance x = x = '@'
let isStoneWall x = x = '#'
let isOpenPassage x = x = '.'
let isKey x = x >= 'a' && x <= 'z'
let isDoor x = x >= 'A' && x <= 'Z'

let map =
    File.ReadAllLines "in"
    |> Seq.mapi (fun y row -> row |> Seq.mapi (fun x col -> (x, y), col))
    |> Seq.collect id
    |> Map.ofSeq

let entrance = map |> Map.findKey (fun _ v -> isEntrance v)

let ensureKey key =
    match isKey key with
    | true -> ()
    | false -> failwith "not a key"

let ensureDoor door =
    match isDoor door with
    | true -> ()
    | false -> failwith "not a door"

type Keychain =
    private
    | Keychain of uint32

    static member empty = Keychain 0u

    member this.value =
        match this with
        | Keychain x -> x

    member this.collect key =
        ensureKey key
        let b = 1u <<< (int key - int 'a')
        Keychain(this.value ||| b)

    member this.remove key =
        ensureKey key
        let b = 1u <<< (int key - int 'a')
        Keychain(this.value ^^^ b)

    member this.missingDoorKey door =
        ensureDoor door
        let b = 1u <<< (int door - int 'A')
        this.value &&& b <> b

let keys = map |> Map.filter (fun _ v -> isKey v)
let all = keys |> Map.fold (fun (s: Keychain) k v -> s.collect v) Keychain.empty

let search (map: Grid) (start: Point) (collected: Keychain) =
    let mutable Q = Set.ofList [ (0, start, collected) ]
    let mutable V = Set.empty
    let mutable S = None

    let directions = [ (0, -1); (1, 0); (0, 1); (-1, 0) ]

    while not Q.IsEmpty && S.IsNone do
        let v = Q.MinimumElement
        Q <- Set.remove v Q

        let steps, pos, collected = v
        let px, py = pos

        if collected = all then
            S <- Some steps
        else
            V <- V |> Set.add (pos, collected)

            for dx, dy in directions do
                let pos' = (px + dx, py + dy)

                match pos' with
                | x when V.Contains(x, collected) -> ()
                | x when isStoneWall map[x] -> ()
                | x when isDoor map[x] && collected.missingDoorKey map[x] -> ()
                | pos' ->
                    let mutable keys' =
                        match isKey map[pos'] with
                        | true -> collected.collect map[pos']
                        | false -> collected

                    Q <- Q |> Set.add (steps + 1, pos', keys')

    match S with
    | Some s -> s
    | None -> failwith "no solution"

let part1 map = search map entrance Keychain.empty
printfn "ans1 %A" (part1 map)

let part2 map =
    let ex, ey = entrance

    let mutable map = map

    [ (0, 0); (0, -1); (-1, 0); (1, 0); (0, 1) ]
    |> Seq.map (fun (dx, dy) -> (ex + dx, ey + dy))
    |> Seq.iter (fun p -> map <- map |> Map.add p '#')

    // top-left quadrant

    let mutable collected = all

    for x in 0 .. ex - 1 do
        for y in 0 .. ey - 1 do
            match isKey map[(x, y)] with
            | true -> collected <- collected.remove map[(x, y)]
            | false -> ()

    let stepsTL = search map (ex - 1, ey - 1) collected

    // top-right quadrant

    let mutable collected = all

    for x in ex + 1 .. ex * 2 do
        for y in 0 .. ey - 1 do
            match isKey map[(x, y)] with
            | true -> collected <- collected.remove map[(x, y)]
            | false -> ()

    let stepsTR = search map (ex + 1, ey - 1) collected

    // bottom-right quadrant

    let mutable collected = all

    for x in ex + 1 .. ex * 2 do
        for y in ey + 1 .. ey * 2 do
            match isKey map[(x, y)] with
            | true -> collected <- collected.remove map[(x, y)]
            | false -> ()

    let stepsBR = search map (ex + 1, ey + 1) collected

    // bottom-left quadrant

    let mutable collected = all

    for x in 0 .. ex - 1 do
        for y in ey + 1 .. ey * 2 do
            match isKey map[(x, y)] with
            | true -> collected <- collected.remove map[(x, y)]
            | false -> ()

    let stepsBL = search map (ex - 1, ey + 1) collected

    stepsTL + stepsTR + stepsBR + stepsBL

printfn "ans2 %A" (part2 map)
