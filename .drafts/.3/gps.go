// =============================================================================
// Auth: Alex Celani
// File: gsp.go
// Revn: 07-13-2023  0.4
// Func: asynchronously receive messages from another program, print
//
// TODO: create
// =============================================================================
// CHANGE LOG
// -----------------------------------------------------------------------------
// 05-04-2022: init
// 05-09-2022: began commenting
// 05-10-2022: finished commenting
//             added alter() and forward()
//             made passed value to alter()/forward() a slice
//             changed refs to recv byte array to slice of n bytes
// 07-13-2023: copied over, removed call to alter()/forward()
//             removed any call to write
//             put read and print in a for loop
//             wrote getFileName() to create a filename based off
//              current time
//             write output to a logfile
//             moved file writing to AFTER checking for fin
//
// =============================================================================

package main

import ( 
    "net"       // ResolveTCPAddr, conn.Write,Read,Close
                // ListenTCP, listener.Accept
    "os"        // Args, Stderr, Exit, Create, Write
    "fmt"       // Fprintf, Println
    "strings"   // ToLower, Compare
    "time"      // Now, Year, Month, Day, Hour, Minute, Second
)


// quick error checking
func check( err error ) {
    // if user checks an error, it better be nil
    if err != nil {
        // put error in string, print in stderr
        fmt.Fprintf( os.Stderr, "Fatal error: %s", err.Error() )
        os.Exit( 2 )        // quit
    }
}


func getFileName() string{
    // capture current time
    t := time.Now()
    // format as YYYYMMDDHHmmss.gps
    return t.Format( "20060102150405" ) + ".gps"
    //return name + // return name
}


// function to handle incoming connections
func handleClient( conn net.Conn ) {
    defer conn.Close()  // barring any error, still close connection
//    defer file.Close()  // barring any error, still close connection

    var buf [512]byte   // declare large byte array, store messages

    // iterate forever to always read over connection
    for {
        //defer file.Close()
        // read n bytes from connection into buffer
        n, err := conn.Read( buf[0:] )
        if err != nil { // erroring on read will simply leave the 
            return      // function so it can start again later
        }

        // print recv'd message
        // string() only works on byte SLICES so [:] is required
        fmt.Println( "recv: ", string( buf[:n] ) )

        // if recv'd message is start, create and open file
        if strings.Compare( string( buf[:n] ), "start" ) == 0 {
            filename := getFileName()   // get filename
            file, err = os.Create( filename )  // create new file
            check( err )    // make sure it worked
            defer file.Close()
            open = true
            break
        }

        if !open {
            break
        }

        // if recv'd message is fin, break forloop
        if strings.Compare( string( buf[:n] ), "fin" ) == 0 {
            file.Close()
            os.Exit( 2 )        // quit
        }

        file, err = os.OpenFile( filename, os.O_WRONLY, 0644 )
        check( err )
        // write info to file
        _, err = file.Write( buf[:n] )
        check( err )    // make sure write worked
        // write newline
        _, err = file.WriteString( "\n" )
        check( err )    // make sure write worked

        file.Sync()

        file.Close()
        /*
        // write that response back to original client
        _, err = conn.Write( []byte( send ) )
        if err != nil { // erroring on write will simply leave the
            return      // function so it can start again later
        }
        */

        //fmt.Println( "sent: ", send )   // print response
    }
}


var file *os.File
var filename string
var open bool


func main() {

    // ip:port
    // ip doesn't exist, implies localhost
    service := ":1202"      // capture ip address and host
    
    // "resolve" ip & host according to TCP rules
    tcpAddr, err := net.ResolveTCPAddr( "tcp", service )
    check( err )    // check error

    // bind and "listen" to ip and port, according to tcp rules
    listener, err := net.ListenTCP( "tcp", tcpAddr )
    check( err )    // check error


    // iterate forever
    // TODO i mean i can totally make this more user friendly
    for {
        // accept a connection that makes its way to bound port
        conn, err := listener.Accept()
        if err != nil {     // if connection fails...
            continue        // don't quit program, not fatal error
        }

        // asynchronous function to handle connection to client
        go handleClient( conn )
    }
}



