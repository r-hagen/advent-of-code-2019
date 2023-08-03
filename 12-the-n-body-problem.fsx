open System
open System.IO
open System.Text.RegularExpressions
open Microsoft.FSharp.Core.Operators.Checked

type Vec = { x: int; y: int; z: int }

type Moon = { pos: Vec; vel: Vec }

let moons =
    File.ReadAllLines "in"
    |> Seq.map (fun line ->
        let pattern = "[xyz]=(-?\d+)"

        let pos =
            Regex.Matches(line, pattern)
            |> Seq.map (fun m -> int m.Groups[1].Value)
            |> Seq.toArray

        { pos = { x = pos[0]; y = pos[1]; z = pos[2] }
          vel = { x = 0; y = 0; z = 0 } })
    |> Seq.toList

let applyGravity moon others =
    let changeVelocity a b =
        if a < b then 1
        else if a > b then -1
        else 0

    let velocity =
        others
        |> Seq.fold
            (fun acc other ->
                { acc with
                    x = acc.x + changeVelocity moon.pos.x other.pos.x
                    y = acc.y + changeVelocity moon.pos.y other.pos.y
                    z = acc.z + changeVelocity moon.pos.z other.pos.z })
            moon.vel

    { moon with vel = velocity }

let applyVelocity moon =
    { moon with
        pos =
            { x = moon.pos.x + moon.vel.x
              y = moon.pos.y + moon.vel.y
              z = moon.pos.z + moon.vel.z } }

let totalEnergy moon =
    let energy v =
        Math.Abs(v.x) + Math.Abs(v.y) + Math.Abs(v.z)

    let pot = energy moon.pos
    let kin = energy moon.vel
    pot * kin

let ans1 =
    let mutable moons = moons |> Seq.toList

    for x in 1..1000 do
        moons <-
            moons
            |> Seq.map (fun m -> applyGravity m moons)
            |> Seq.map (applyVelocity)
            |> Seq.toList

        let te = moons |> Seq.sumBy (totalEnergy)
        printfn $"Iteration {x}: {te}"

let rec gcd a b =
    match a, b with
    | (a, 0L) -> a
    | (a, b) -> gcd b (a % b)

let lcm a b = a * b / gcd a b

let ans2 =
    let mx = moons |> Seq.map (fun m -> m.pos.x, m.vel.x) |> Set.ofSeq
    let my = moons |> Seq.map (fun m -> m.pos.y, m.vel.y) |> Set.ofSeq
    let mz = moons |> Seq.map (fun m -> m.pos.z, m.vel.z) |> Set.ofSeq

    let mutable moons2 = moons |> Seq.toList

    let mutable step_number = 1L
    let mutable steps_x = -1L
    let mutable steps_y = -1L
    let mutable steps_z = -1L

    while steps_x = -1 || steps_y = -1 || steps_z = -1 do
        moons2 <-
            moons2
            |> Seq.map (fun m -> applyGravity m moons2)
            |> Seq.map (applyVelocity)
            |> Seq.toList

        let sx = moons2 |> Seq.map (fun m -> m.pos.x, m.vel.x) |> Set.ofSeq
        let sy = moons2 |> Seq.map (fun m -> m.pos.y, m.vel.y) |> Set.ofSeq
        let sz = moons2 |> Seq.map (fun m -> m.pos.z, m.vel.z) |> Set.ofSeq

        if mx = sx && steps_x = -1 then
            steps_x <- step_number

        if my = sy && steps_y = -1 then
            steps_y <- step_number

        if mz = sz && steps_z = -1 then
            steps_z <- step_number

        step_number <- step_number + 1L

    printfn "ans2: %A" (lcm (lcm steps_x steps_y) steps_z)
