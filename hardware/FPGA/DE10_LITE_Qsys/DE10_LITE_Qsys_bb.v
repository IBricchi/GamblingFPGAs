
module DE10_LITE_Qsys (
	accelerometer_spi_external_interface_I2C_SDAT,
	accelerometer_spi_external_interface_I2C_SCLK,
	accelerometer_spi_external_interface_G_SENSOR_CS_N,
	accelerometer_spi_external_interface_G_SENSOR_INT,
	altpll_0_areset_conduit_export,
	altpll_0_locked_conduit_export,
	bet1_0_conduit_end_datac,
	button_external_connection_export,
	clk_clk,
	clk_sdram_clk,
	hex_0_external_connection_export,
	hex_1_external_connection_export,
	hex_2_external_connection_export,
	hex_3_external_connection_export,
	hex_4_external_connection_export,
	hex_5_external_connection_export,
	key_external_connection_export,
	led_external_connection_export,
	reset_reset_n,
	sdram_wire_addr,
	sdram_wire_ba,
	sdram_wire_cas_n,
	sdram_wire_cke,
	sdram_wire_cs_n,
	sdram_wire_dq,
	sdram_wire_dqm,
	sdram_wire_ras_n,
	sdram_wire_we_n,
	switch_external_connection_export,
	tilt3_0_conduit_end_datac,
	tilt4_0_conduit_end_datac);	

	inout		accelerometer_spi_external_interface_I2C_SDAT;
	output		accelerometer_spi_external_interface_I2C_SCLK;
	output		accelerometer_spi_external_interface_G_SENSOR_CS_N;
	input		accelerometer_spi_external_interface_G_SENSOR_INT;
	input		altpll_0_areset_conduit_export;
	output		altpll_0_locked_conduit_export;
	input	[5:0]	bet1_0_conduit_end_datac;
	input	[3:0]	button_external_connection_export;
	input		clk_clk;
	output		clk_sdram_clk;
	output	[6:0]	hex_0_external_connection_export;
	output	[6:0]	hex_1_external_connection_export;
	output	[6:0]	hex_2_external_connection_export;
	output	[6:0]	hex_3_external_connection_export;
	output	[6:0]	hex_4_external_connection_export;
	output	[6:0]	hex_5_external_connection_export;
	input	[3:0]	key_external_connection_export;
	output	[9:0]	led_external_connection_export;
	input		reset_reset_n;
	output	[12:0]	sdram_wire_addr;
	output	[1:0]	sdram_wire_ba;
	output		sdram_wire_cas_n;
	output		sdram_wire_cke;
	output		sdram_wire_cs_n;
	inout	[15:0]	sdram_wire_dq;
	output	[1:0]	sdram_wire_dqm;
	output		sdram_wire_ras_n;
	output		sdram_wire_we_n;
	input	[9:0]	switch_external_connection_export;
	input		tilt3_0_conduit_end_datac;
	input		tilt4_0_conduit_end_datac;
endmodule
