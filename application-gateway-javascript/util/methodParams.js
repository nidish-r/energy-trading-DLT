// Path to org1 user private key directory.
const swappingNetwork1Init = {
  methodName: 'InitializeSSNetwork',
  methodParams: ['Network-1',
  'TruePower',
  'Active'
]}

const swappingNetwork2Init = {
  methodName: 'InitializeSSNetwork',
  methodParams: ['Network-2',
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
    'Jio'
]}

const swappingStation2Init = {
  methodName: 'InitializeSwappingStation',
  methodParams: [
    'SS-2',
    'Delta 1',
    'Network-2',
    '12.990363054, 77.5884480123',
    '39/14, Sarjapur - Marathahalli Rd, Ibbaluru, Bellandur, Bengaluru, Karnataka 560103',
    'Delta-3000',
    'business@jio.com',
    '1800 889 9999',
    'Delta'
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
    'Network-2',
    '010522',
    '1676456136',
]}

const userInit = {
  methodName: 'InitializeUser',
  methodParams: [
    'User-1',
    'bkcninja',
    'Sarjapur Road, Bangalore Bangalore East, Pin Code: 560035',
    '0000 0000 0000',
    'bkcninja@dltmail.com',
    '9989989989'
]}

const rechargeUserWallet = {
  methodName: 'RechargeUserWallet',
  methodParams: [
    'User-1',
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
    'SS-2',
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

const transferFromSS1toUser = {
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

const transferFromSS2toUser = {
  methodName: 'TransferBatteryFromSSToUser',
  methodParams: [
    'Battery-2',
    '100',
    '100',
    '40',
    '4000',
    'User-1',
    '500'
]}

const transferFromUsertoSS1 = {
  methodName: 'TransferBatteryFromUserToSS',
  methodParams: [
    'Battery-2',
    '10',
    '100',
    '4',
    '4000',
    'SS-1',
]}

const transferFromUsertoSS2 = {
  methodName: 'TransferBatteryFromUserToSS',
  methodParams: [
    'Battery-1',
    '10',
    '100',
    '4',
    '4000',
    'SS-2',
]}

const readSSNetwork1 = {
  methodName: 'ReadSSNetwork',
  methodParams: [
    'Network-1',
]}

const readSSNetwork2 = {
  methodName: 'ReadSSNetwork',
  methodParams: [
    'Network-2',
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

const readUser = {
  methodName: 'ReadUser',
  methodParams: [
    'User-1',
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
  1: [swappingNetwork1Init, swappingNetwork2Init],
  2: [swappingStation1Init, swappingStation2Init],
  3: [batteryInit1, batteryInit2],
  4: [userInit], 
  5: [rechargeUserWallet],
  6: [dockBattery1, dockBattery2],
  7: [returnFromService1, returnFromService2],
  8: [transferFromSS1toUser],
  9: [transferFromUsertoSS2],
  10: [transferFromSS2toUser],
  11: [transferFromUsertoSS1]
}

module.exports = { parameterMap, 
  readBattery1, 
  readBattery2, 
  readSSNetwork1, 
  readSSNetwork2, 
  readSwappingStation1,
  readSwappingStation2,
  readUser };