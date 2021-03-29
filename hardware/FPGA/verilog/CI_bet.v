module CI_bet
	#(  parameter LATENCY=0  )
	(
		input 		reset,
		input		clk,
		input		clk_en,
		input	[31:0]	dataa,
		input	[31:0]	datab,
		input	[ 5:0]	datac,
		input		start,
		output	[31:0]	result,
		output		done
	);

	wire [31:0] sqrt_result;

	bet Inst_bet(
		.aclr	(reset),
		.clk_en	(clk_en),
		.clock	(clk),
		.dataa	(dataa),
		.datab	(datab),
		.datac	(datac),
		.result	(sqrt_result)
	);

	
	// Statemachine timing the done bit 
	reg [4:0] state_sync;

	always @ (posedge clk) begin
		if (reset) begin
			state_sync <= 0;
		end else if (start & clk_en) begin
			state_sync <= LATENCY;			// Floating point unit is starting 
		end else if (clk_en) begin
			state_sync <= state_sync - 1;		// Floating point isn't done yet
		end else begin
			state_sync <= 0;
		end
	end

	assign done = (state_sync == 1) & clk_en;
	assign result = sqrt_result;
endmodule

