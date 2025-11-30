package main

/*
Design an Uber-like service managing drivers, riders, and trip assignments.
-> user(driver, rider)
	as rider -> enter starting and destination point see the fare for different vehicle types
			-> rider can confirm the booking for a vehicle type
			-> a driver should be assigned for that ride
			-> rider should see the driver location in real time.
	as driver -> driver should accept the ride.
			  -> driver should see the pickup and drop point for a ride and the final fare
			  -> driver should verify the rider with otp authentication

Driver
Rider
Fare
commuteTime
Ride
Vehicle
commuteTimeAlgo
FareAlgo

GetRide()
GetFare()
GetCommuteTime()
ConfirmRide()
AcceptRide()
NotifyDriver()
GetDriverLocation()
SendLocation()
NearPickupPoint()
StartRide()
EndRide()

FareService
	-> GetFare()
CommuteService
	-> CommuteTime(vehicleType) ->

RideService
	-> GetRide()
	-> StartRide()
	-> EndRide()
	-> AcceptRide()

DriverService
	-> SendLocation() -> ns.NotifyRider
	-> GetLocation()

RiderService
	-> GetLocation()
	-> SendLocation()

NotificationService
	-> NotifyRider
	-> NotifyDriver


*/

