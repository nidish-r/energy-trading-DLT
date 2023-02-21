const SimulationMap = {
  0 : "BEGIN SIMULATION",
  1: "Swapping Networks Init",
  2: "Swapping Stations Init",
  3: "Batteries Init",
  4: "User Init",
  5: "User Recharge",
  6: "Dock Batteries",
  7: "Return From Service",
  8: "Transfer Battery From SS1 To User",
  9: "Transfer Battery From User To SS2",
  10: "Transfer Battery From SS2 To User",
  11: "Transfer Battery From User To SS1",
  12: "SIMULATION FINISHED"
}

function progressHelper(simulationPhase) {
  const phaseArray = []
  const steps = []
  for(let i = 0; i <= simulationPhase + 1; i++) { 
    if(i < simulationPhase + 1) {
      phaseArray.push(SimulationMap[i]);
    } else {
      if(SimulationMap[i] != undefined ) {
        phaseArray.push("NEXT: " + SimulationMap[i]);
      }
    }
  }

  for(let i = 0; i <= 12; i++) { 
    steps.push(SimulationMap[i]);
  }

  return { phaseArray, steps };
}

module.exports = { SimulationMap, progressHelper }