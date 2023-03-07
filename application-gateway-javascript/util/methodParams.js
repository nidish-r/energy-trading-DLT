// Path to org1 user private key directory.
const swappingNetwork1Init = {
  methodName: 'InitializeSSNetwork',
  methodParams: ['Network-1',
  'TruePower',
  'Active'
]}

const swappingNetwork2Init = {
  methodName: 'InitializeSSNetwork',
  methodParams: ['Network-1',
  'Delta',
  'Active'
]}

const swappingStation1Init = {
  methodName: 'InitializeSwappingStation',
  methodParams: [
    'SS-1',
    'Jio AVANA 1',
    'Network-1',
    '12.990363054, 77.5884480123',
    '39/14, Sarjapur - Marathahalli Rd, Ibbaluru, Bellandur, Bengaluru, Karnataka 560103',
    'Jio-3000',
    'business@jio.com',
    '1800 889 9999',
    'Zomato',
]}

const swappingStation2Init = {
  methodName: 'InitializeSwappingStation',
  methodParams: [
    'SS-2',
    'Delta 1',
    'Network-1',
    '12.990363054, 77.5884480123',
    '39/14, Sarjapur - Marathahalli Rd, Ibbaluru, Bellandur, Bengaluru, Karnataka 560103',
    'Delta-3000',
    'business@jio.com',
    '1800 889 9999',
    'Zomato',
]}

const fleetInit = {
  methodName: 'InitializeFleet',
  methodParams: ['Fleet-1',
  'West Zone',
  'Zomato',
  'Food Delivery',
  'info@zomato.com',
  '1800 889 9999',
  '39/14, Sarjapur - Marathahalli Rd, Ibbaluru, Bellandur, Bengaluru, Karnataka 560103',
]}

const batteryInit1 = {
  methodName: 'InitializeBattery',
  methodParams: [
    'Battery-1',
    'EV-12V-80AH',
    '0',
    '100',
    '0',
    '4000',
    'Network-1',
    '010521',
    '1676456136',
]}


const batteryInit2 = {
  methodName: 'InitializeBattery',
  methodParams: [
    'Battery-2',
    'EV-12V-80AH',
    '0',
    '100',
    '0',
    '4000',
    'Network-1',
    '010522',
    '1676456136',
]}

const batteryInit3 = {
  methodName: 'InitializeBattery',
  methodParams: [
    'Battery-3',
    'EV-12V-80AH',
    '0',
    '100',
    '0',
    '4000',
    'Network-1',
    '010522',
    '1676456136',
]}

const batteryInit4 = {
  methodName: 'InitializeBattery',
  methodParams: [
    'Battery-4',
    'EV-12V-80AH',
    '0',
    '100',
    '0',
    '4000',
    'Network-1',
    '010522',
    '1676456136',
]}

const user1Init = {
  methodName: 'InitializeUser',
  methodParams: [
    'User-1',
    'bkcninja',
    'Sarjapur Road, Bangalore Bangalore East, Pin Code: 560035',
    '0000 0000 0000',
    'bkcninja@dltmail.com',
    'Fleet-1',
    'Zomato',
    '9989989989'
]}

const user2Init = {
  methodName: 'InitializeUser',
  methodParams: [
    'User-2',
    'hghninja',
    'Sarjapur Road, Bangalore Bangalore East, Pin Code: 560035',
    '0000 0000 0000',
    'hghinja@dltmail.com',
    'Fleet-1',
    'Zomato',
    '9989989999'
]}


const rechargeUser1Wallet = {
  methodName: 'RechargeUserWallet',
  methodParams: [
    'User-1',
    '1000'
]}

const rechargeUser2Wallet = {
  methodName: 'RechargeUserWallet',
  methodParams: [
    'User-2',
    '1000'
]}

const dockBattery1 = {
  methodName: 'DockBatteryOnSwappingStation',
  methodParams: [
    'Battery-1',
    '0',
    '100',
    '0',
    '4000',
    'SS-1',
]}

const dockBattery2 = {
  methodName: 'DockBatteryOnSwappingStation',
  methodParams: [
    'Battery-2',
    '0',
    '100',
    '0',
    '4000',
    'SS-1',
]}

const dockBattery3 = {
  methodName: 'DockBatteryOnSwappingStation',
  methodParams: [
    'Battery-3',
    '0',
    '100',
    '0',
    '4000',
    'SS-2',
]}

const dockBattery4 = {
  methodName: 'DockBatteryOnSwappingStation',
  methodParams: [
    'Battery-4',
    '0',
    '100',
    '0',
    '4000',
    'SS-2',
]}


const allocateBattery1ToFleet = {
  methodName: 'AllocateBatteryToFleet',
  methodParams: [
    'Fleet-1',
    'Battery-1',
    'Undocked',
    'Zomato',
]}

const allocateBattery2ToFleet = {
  methodName: 'AllocateBatteryToFleet',
  methodParams: [
    'Fleet-1',
    'Battery-2',
    'Undocked',
    'Zomato',
]}

const allocateBattery3ToFleet = {
  methodName: 'AllocateBatteryToFleet',
  methodParams: [
    'Fleet-1',
    'Battery-3',
    'Undocked',
    'Zomato',
]}

const allocateBattery4ToFleet = {
  methodName: 'AllocateBatteryToFleet',
  methodParams: [
    'Fleet-1',
    'Battery-4',
    'Undocked',
    'Zomato',
]}

const returnFromService1 = {
  methodName: 'ReturnBatteryFromService',
  methodParams: [
    'Battery-1',
    '100',
    '100',
    '40',
    '4000',
]}

const returnFromService2 = {
  methodName: 'ReturnBatteryFromService',
  methodParams: [
    'Battery-2',
    '100',
    '100',
    '40',
    '4000',
]}

const returnFromService3 = {
  methodName: 'ReturnBatteryFromService',
  methodParams: [
    'Battery-3',
    '100',
    '100',
    '40',
    '4000',
]}

const returnFromService4 = {
  methodName: 'ReturnBatteryFromService',
  methodParams: [
    'Battery-4',
    '100',
    '100',
    '40',
    '4000',
]}

// const DockOONBatteryOnSwappingStation = {
//   methodName: 'DockOONBatteryOnSwappingStation',
//   methodParams: [
//     'Battery-U1',
//     'Network-1',
//     'SS-1',
//     'User-1',
//     'EV-12V-80AH',
//     '010521',
//     '1676456136',
// ]}

// const VerifiyOONBatteryOnSS = {
//   methodName: 'VerifiyOONBatteryOnSS',
//   methodParams: [
//     'Battery-U1',
//     '0',
//     '80',
//     '0',
//     '2000',
// ]}

const transferFromSS1toUser1 = {
  methodName: 'TransferBatteryFromSSToUser',
  methodParams: [
    'Battery-1',
    '100',
    '100',
    '40',
    '4000',
    'User-1',
    '500'
]}

const transferFromSS2toUser1 = {
  methodName: 'TransferBatteryFromSSToUser',
  methodParams: [
    'Battery-4',
    '100',
    '100',
    '40',
    '4000',
    'User-1',
    '500'
]}

const transferFromSS2toUser2 = {
  methodName: 'TransferBatteryFromSSToUser',
  methodParams: [
    'Battery-3',
    '100',
    '100',
    '40',
    '4000',
    'User-2',
    '500'
]}

const transferFromSS1toUser2 = {
  methodName: 'TransferBatteryFromSSToUser',
  methodParams: [
    'Battery-2',
    '100',
    '100',
    '40',
    '4000',
    'User-2',
    '500'
]}

const transferFromUser1toSS1 = {
  methodName: 'TransferBatteryFromUserToSS',
  methodParams: [
    'Battery-4',
    '10',
    '100',
    '4',
    '4000',
    'SS-1',
]}

const transferFromUser1toSS2 = {
  methodName: 'TransferBatteryFromUserToSS',
  methodParams: [
    'Battery-1',
    '10',
    '100',
    '4',
    '4000',
    'SS-2',
]}

const transferFromUser2toSS2 = {
  methodName: 'TransferBatteryFromUserToSS',
  methodParams: [
    'Battery-2',
    '10',
    '100',
    '4',
    '4000',
    'SS-2',
]}

const transferFromUser2toSS1 = {
  methodName: 'TransferBatteryFromUserToSS',
  methodParams: [
    'Battery-3',
    '10',
    '100',
    '4',
    '4000',
    'SS-1',
]}


const readSSNetwork1 = {
  methodName: 'ReadSSNetwork',
  methodParams: [
    'Network-1',
]}

const readSSNetwork2 = {
  methodName: 'ReadSSNetwork',
  methodParams: [
    'Network-1',
]}

const readSwappingStation1 = {
  methodName: 'ReadSwappingStation',
  methodParams: [
    'SS-1',
]}

const readSwappingStation2 = {
  methodName: 'ReadSwappingStation',
  methodParams: [
    'SS-2',
]}

const readUser1 = {
  methodName: 'ReadUser',
  methodParams: [
    'User-1',
]}

const readUser2 = {
  methodName: 'ReadUser',
  methodParams: [
    'User-2',
]}

const readBattery1 = {
  methodName: 'ReadBattery',
  methodParams: [
    'Battery-1',
]}

const readBattery2 = {
  methodName: 'ReadBattery',
  methodParams: [
    'Battery-2',
]}

const readBattery3 = {
  methodName: 'ReadBattery',
  methodParams: [
    'Battery-3',
]}

const readBattery4 = {
  methodName: 'ReadBattery',
  methodParams: [
    'Battery-4',
]}

// const readUserBattery = {
//   methodName: 'ReadBattery',
//   methodParams: [
//     'Battery-U1',
// ]}


const readBatteryHistory1 = {
  methodName: 'ReadBatteryHistory',
  methodParams: [
    'Battery-1',
]}

const readBatteryHistory2 = {
  methodName: 'ReadBatteryHistory',
  methodParams: [
    'Battery-2',
]}

const parameterMap = {
  1: [swappingNetwork1Init],
  2: [swappingStation1Init, swappingStation2Init],
  3: [batteryInit1, batteryInit2, batteryInit3, batteryInit4],
  4: [fleetInit],
  5: [user1Init, user2Init], 
  6: [rechargeUser1Wallet, rechargeUser2Wallet],
  7: [allocateBattery1ToFleet, allocateBattery2ToFleet, allocateBattery3ToFleet, allocateBattery4ToFleet],
  8: [dockBattery1, dockBattery2, dockBattery3, dockBattery4],
  9: [returnFromService1, returnFromService2, returnFromService3, returnFromService4],
  10: [transferFromSS1toUser1, transferFromSS2toUser2],
  11: [transferFromUser1toSS2, transferFromUser2toSS1],
  12: [transferFromSS2toUser1, transferFromSS1toUser2],
  13: [transferFromUser1toSS1, transferFromUser2toSS2]
}

module.exports = { parameterMap, 
  readBattery1, 
  readBattery2, 
  readBattery3,
  readBattery4,
  readSSNetwork1,
  readSwappingStation1,
  readSwappingStation2,
  readUser1,
  readUser2
};