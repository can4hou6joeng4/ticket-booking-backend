start:
	@docker compose -p ticket-booking up --build

stop:
	@docker compose -p ticket-booking rm -v --force --stop
	@docker rmi ticket-booking
