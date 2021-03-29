// synopsys translate_off
'timescale 1 ps / 1 ps
// synopsys translate_on
module fir (
	aclr,
	clk_en,
	clock,
	dataa,
	datab,
	result
);

	input	aclr;
	input	clk_en;
	input	clock;
	input	[31:0]	dataa;
	input	[31:0]	datab;
	output	[31:0]	result;

	reg 	[31:0] 	result; // = dataa[15:0] << datab[15:0];
	
	reg	 	newClk;
	reg 	[14:0] 	count; 

    reg [31:0] dataInEx;
	reg [31:0] stage1;
	reg [31:0] stage2;
	reg [31:0] stage3;
	reg [31:0] stage4;
	reg [31:0] stage5;
	reg [31:0] stage6;
	reg [31:0] stage7;
	reg [31:0] stage8;
	reg [31:0] stage9;
	reg [31:0] stage10;
	reg [31:0] stage11;
	reg [31:0] stage12;
	reg [31:0] stage13;
	reg [31:0] stage14;
	reg [31:0] stage15;
	reg [31:0] stage16;
	reg [31:0] stage17;
	reg [31:0] stage18;
	reg [31:0] stage19;
	reg [31:0] stage20;
	reg [31:0] stage21;
	reg [31:0] stage22;
	reg [31:0] stage23;
	reg [31:0] stage24;
	reg [31:0] stage25;
	reg [31:0] stage26;
	reg [31:0] stage27;
	reg [31:0] stage28;
	reg [31:0] stage29;
	reg [31:0] stage30;
	reg [31:0] stage31;
	reg [31:0] stage32;
	reg [31:0] stage33;
	reg [31:0] stage34;
	reg [31:0] stage35;
	reg [31:0] stage36;
	reg [31:0] stage37;
	reg [31:0] stage38;
	reg [31:0] stage39;
	reg [31:0] stage40;
	reg [31:0] stage41;
	reg [31:0] stage42;
	reg [31:0] stage43;
	reg [31:0] stage44;
	reg [31:0] stage45;
	reg [31:0] stage46;
	reg [31:0] stage47;
	reg [31:0] stage48;
	reg [31:0] stage49;
	reg [31:0] intermed;


	always @ (dataa) begin
		
		stage1 	<= dataa;
		stage2 	<= stage1;
		stage3 	<= stage2;
		stage4 	<= stage3;
		stage5 	<= stage4;
		stage6 	<= stage5;
		stage7 	<= stage6;
		stage8 	<= stage7;
		stage9 	<= stage8;
		stage10	<= stage9;
		stage11	<= stage10;
		stage12 <= stage11;
		stage13 <= stage12;
		stage14	<= stage13;
		stage15	<= stage14;
		stage16	<= stage15;
		stage17	<= stage16;
		stage18	<= stage17;
		stage19	<= stage18;
		stage20 <= stage19;
		stage21 <= stage20;
		stage22 <= stage21;
		stage23 <= stage22;
		stage24 <= stage23;
		stage25 <= stage24;
		stage26 <= stage25;
		stage27 <= stage26;
		stage28 <= stage27;
		stage29 <= stage28;
		stage30 <= stage29;
		stage31 <= stage30;
		stage32 <= stage31;
		stage33 <= stage32;
		stage34 <= stage33;
		stage35 <= stage34;
		stage36 <= stage35;
		stage37 <= stage36;
		stage38 <= stage37;
		stage39 <= stage38;
		stage40 <= stage39;
		stage41 <= stage40;
		stage42 <= stage41;
		stage43 <= stage42;
		stage44 <= stage43;
		stage45 <= stage44;
		stage46 <= stage45;
		stage47 <= stage46;
		stage48 <= stage47;
		stage49 <= stage48;	
		
intermed <= (1*dataInEx) + (2*stage1) + (-1*stage2) + (-2*stage3) + (1*stage4) + (0*stage5) + (-2*stage6) + (1*stage7) + (1*stage8) + (-3*stage9) + (1*stage10) +  (3*stage11) + (-5*stage12) + (0*stage13) + (6*stage14) + (-6*stage15) + (-3*stage16) + (10*stage17) + (-7*stage18) + (-9*stage19) + (19*stage20) + (-7*stage21) + (-31*stage22) + (74*stage23) + (163*stage24) + (74*stage25) + (-31*stage26) + (-7*stage27) + (19*stage28) + (-9*stage29) + (-7*stage30) + (10*stage31) + (-3*stage32) + (-6*stage33) + (6*stage34)  + (0*stage35) + (-5*stage36) + (3*stage37) + (1*stage38) + (-3*stage39) + (1*stage40)  + (1*stage41) + (-2*stage42) + (0*stage43) + (1*stage44)  + (-2*stage45) + (-1*stage46) + (2*stage47) + (1*stage48) + (1*stage49);

result <= intermed >> 8;
		
	end

endmodule
